import http from 'k6/http';
import { check, fail, group, sleep } from 'k6';
import { waitForReadiness } from './helpers/readiness.js';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');

http.setResponseCallback(
  http.expectedStatuses({ min: 200, max: 399 }, 400)
);

export const options = {
  vus: 1,
  duration: '30s',
  setupTimeout: '2m',
  thresholds: {
    'http_req_failed{traffic_type:workload}': ['rate<0.01'],
    'http_req_duration{traffic_type:workload}': ['p(95)<500'],
    checks: ['rate>0.99'],
  },
};

const workloadParams = {
  tags: {
    traffic_type: 'workload',
  },
};

const jsonRequestParams = {
  headers: {
    'Content-Type': 'application/json',
  },
  tags: {
    traffic_type: 'workload',
  },
};

export function setup() {
  waitForReadiness(BASE_URL);
}

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
    const live = http.get(url('/healthz'), workloadParams);
    const liveBody = parseJson(live);

    const liveOk = check(live, {
      'healthz status is 200': (r) => r.status === 200,
      'healthz body is ok': () => liveBody && liveBody.status === 'ok',
    });

    if (!liveOk) {
      fail('service liveness check failed');
    }

    const ready = http.get(url('/readyz'), workloadParams);
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
    const products = http.get(url('/api/v1/product?limit=5'), workloadParams);
    const productsBody = parseJson(products);

    check(products, {
      'product list status is 200': (r) => r.status === 200,
      'product list has success envelope': () => hasSuccessEnvelope(productsBody),
      'product list data has items array': () =>
        productsBody && productsBody.data && Array.isArray(productsBody.data.items),
    });

    const categories = http.get(url('/api/v1/category'), workloadParams);
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
    const res = http.post(url('/api/v1/auth/refresh-token'), JSON.stringify({}), jsonRequestParams);
    const body = parseJson(res);

    check(res, {
      'invalid refresh request returns 400': (r) => r.status === 400,
      'invalid refresh request has error envelope': () =>
        body && body.meta && body.meta.code === 400,
    });
  });

  sleep(1);

  group('Metrics Endpoint', () => {
    const res = http.get(url('/metrics'), workloadParams);

    check(res, {
      'metrics status is 200': (r) => r.status === 200,
      'metrics exposes prometheus text': (r) =>
        r.body.includes('# HELP') && r.body.includes('go_goroutines'),
    });
  });

  sleep(1);
}
