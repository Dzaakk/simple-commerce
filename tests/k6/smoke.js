import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { fail } from 'k6';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');

http.setResponseCallback(
  http.expectedStatuses({ min: 200, max: 399 }, 400)
);

export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500'],
    checks: ['rate>0.99'],
  },
};

const jsonHeaders = {
  headers: {
    'Content-Type': 'application/json',
  },
};

function url(path) {
  return `${BASE_URL}${path}`;
}

function parseJson(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

function hasSuccessEnvelope(body) {
  return body && body.meta && body.meta.code === 200 && body.meta.message === 'Success';
}

export default function () {
  group('Health Checks', () => {
    const live = http.get(url('/healthz'));
    const liveBody = parseJson(live);

    const liveOk = check(live, {
      'healthz status is 200': (r) => r.status === 200,
      'healthz body is ok': () => liveBody && liveBody.status === 'ok',
    });

    if (!liveOk) {
      fail('service liveness check failed');
    }

    const ready = http.get(url('/readyz'));
    const readyBody = parseJson(ready);

    const readyOk = check(ready, {
      'readyz status is 200': (r) => r.status === 200,
      'readyz reports postgres ok': () =>
        readyBody && readyBody.dependencies && readyBody.dependencies.postgres === 'ok',
      'readyz reports redis ok': () =>
        readyBody && readyBody.dependencies && readyBody.dependencies.redis === 'ok',
    });

    if (!readyOk) {
      fail('service readiness check failed');
    }
  });

  sleep(1);

  group('Catalog Read Endpoints', () => {
    const products = http.get(url('/api/v1/product?limit=5'));
    const productsBody = parseJson(products);

    check(products, {
      'product list status is 200': (r) => r.status === 200,
      'product list has success envelope': () => hasSuccessEnvelope(productsBody),
      'product list data has items array': () =>
        productsBody && productsBody.data && Array.isArray(productsBody.data.items),
    });

    const categories = http.get(url('/api/v1/category'));
    const categoriesBody = parseJson(categories);

    check(categories, {
      'category list status is 200': (r) => r.status === 200,
      'category list has success envelope': () => hasSuccessEnvelope(categoriesBody),
      'category list data is array': () =>
        categoriesBody && Array.isArray(categoriesBody.data),
    });
  });

  sleep(1);

  group('Safe Write Contract', () => {
    const res = http.post(url('/api/v1/auth/refresh-token'), JSON.stringify({}), jsonHeaders);
    const body = parseJson(res);

    check(res, {
      'invalid refresh request returns 400': (r) => r.status === 400,
      'invalid refresh request has error envelope': () =>
        body && body.meta && body.meta.code === 400,
    });
  });

  sleep(1);

  group('Metrics Endpoint', () => {
    const res = http.get(url('/metrics'));

    check(res, {
      'metrics status is 200': (r) => r.status === 200,
      'metrics exposes prometheus text': (r) =>
        r.body.includes('# HELP') && r.body.includes('go_goroutines'),
    });
  });

  sleep(1);
}
