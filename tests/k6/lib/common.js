import http from 'k6/http';
import { check, sleep } from 'k6';
import { fail } from 'k6';

export const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');

http.setResponseCallback(
  http.expectedStatuses({ min: 200, max: 399 }, 400)
);

export function url(path) {
  return `${BASE_URL}${path}`;
}

export function parseJson(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

export function hasSuccessEnvelope(body) {
  return body && body.meta && body.meta.code === 200 && body.meta.message === 'Success';
}

export function assertReady() {
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
    'readyz postgres is ok': () =>
      readyBody && readyBody.dependencies && readyBody.dependencies.postgres === 'ok',
    'readyz redis is ok': () =>
      readyBody && readyBody.dependencies && readyBody.dependencies.redis === 'ok',
  });

  if (!readyOk) {
    fail('service readiness check failed');
  }
}

export function randomItem(items) {
  if (!items || items.length === 0) {
    return null;
  }

  return items[Math.floor(Math.random() * items.length)];
}

export function thinkTime(minSeconds = 0.3, maxSeconds = 1.2) {
  const delay = minSeconds + Math.random() * (maxSeconds - minSeconds);
  sleep(delay);
}

export function encodeQuery(params) {
  const pairs = [];

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') {
      continue;
    }

    pairs.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
  }

  return pairs.length > 0 ? `?${pairs.join('&')}` : '';
}
