import http from 'k6/http';
import { check, fail } from 'k6';

export const options = { vus: 1, iterations: 1 };

const baseURL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const ready = http.get(`${baseURL}/readyz`);
  if (!check(ready, { 'service is ready': (r) => r.status === 200 })) {
    fail(`readiness failed: ${ready.status} ${ready.body}`);
  }

  const v1 = http.get(`${baseURL}/api/v1/product?limit=5&sort_by=newest`);
  const v2 = http.get(`${baseURL}/api/v2/product?limit=5&sort_by=newest`);
  check(v1, { 'v1 list succeeds': (r) => r.status === 200 });
  check(v2, { 'v2 list succeeds': (r) => r.status === 200 });

  const v1Body = v1.json();
  const v2Body = v2.json();
  check(null, {
    'v1 and v2 return identical data': () => JSON.stringify(v1Body.data) === JSON.stringify(v2Body.data),
    'seed returns products': () => v1Body.data.items.length === 5,
  });

  const id = v1Body.data.items[0].id;
  const detailV1 = http.get(`${baseURL}/api/v1/product/${id}`);
  const detailV2 = http.get(`${baseURL}/api/v2/product/${id}`);
  check(null, {
    'v1 detail succeeds': () => detailV1.status === 200,
    'v2 detail succeeds': () => detailV2.status === 200,
    'detail contracts match': () => JSON.stringify(detailV1.json().data) === JSON.stringify(detailV2.json().data),
  });
}
