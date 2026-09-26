import http from 'k6/http';
import { check, sleep } from 'k6';

const baseURL = __ENV.BASE_URL || 'http://localhost:8080';
const endpoint = __ENV.ENDPOINT || '/api/v1/product';

export const options = {
  scenarios: {
    catalog: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 100 },
        { duration: __ENV.DURATION || '3m', target: 100 },
        { duration: '30s', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500'],
  },
};

export default function () {
  const page = http.get(`${baseURL}${endpoint}?limit=20&sort_by=newest`);
  const ok = check(page, {
    'list status is 200': (r) => r.status === 200,
    'list has items': (r) => Array.isArray(r.json('data.items')) && r.json('data.items').length > 0,
  });

  if (ok) {
    const items = page.json('data.items');
    const product = items[Math.floor(Math.random() * items.length)];
    check(http.get(`${baseURL}${endpoint}/${product.id}`), {
      'detail status is 200': (r) => r.status === 200,
    });
  }
  sleep(1);
}
