import http from 'k6/http';
import { check, sleep } from 'k6';
import { productIds, categoryIds } from './fixtures/catalog-fixture.js';
import { waitForReadiness } from './helpers/readiness.js';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');
const PRODUCT_LIST_ENDPOINT = __ENV.PRODUCT_LIST_ENDPOINT || '/api/v1/product';
const PRODUCT_DETAIL_ENDPOINT = __ENV.PRODUCT_DETAIL_ENDPOINT || '/api/v1/product';
const PRODUCT_LIMIT = Number(__ENV.PRODUCT_LIMIT || 100);
const TARGET_VUS = positiveInteger(__ENV.TARGET_VUS, 100, 'TARGET_VUS');
const RAMP_UP_DURATION = __ENV.RAMP_UP_DURATION || '30s';
const STEADY_DURATION = __ENV.STEADY_DURATION || '3m';
const RAMP_DOWN_DURATION = __ENV.RAMP_DOWN_DURATION || '30s';

export const options = {
  setupTimeout: '2m',
  scenarios: {
    catalog_browsing: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: RAMP_UP_DURATION, target: TARGET_VUS },
        { duration: STEADY_DURATION, target: TARGET_VUS },
        { duration: RAMP_DOWN_DURATION, target: 0 },
      ],
      gracefulRampDown: RAMP_DOWN_DURATION,
    },
  },
  thresholds: {
    'http_req_failed{traffic_type:workload}': ['rate<0.01'],
    'http_req_duration{traffic_type:workload}': ['p(95)<2000'],
    checks: ['rate>0.99'],
  },
};

const workloadParams = {
  tags: {
    traffic_type: 'workload',
  },
};

export function setup() {
  waitForReadiness(BASE_URL);
}

function positiveInteger(value, fallback, name) {
  if (value === undefined || value === '') {
    return fallback;
  }

  const parsed = Number(value);
  if (!Number.isInteger(parsed) || parsed <= 0) {
    throw new Error(`${name} must be a positive integer, got: ${value}`);
  }

  return parsed;
}

function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function randomChoice(items) {
  return items[Math.floor(Math.random() * items.length)];
}

function checkResponse(res, name) {
  check(res, {
    [`${name} status is 200`]: (r) => r.status === 200,
    [`${name} response time < 2s`]: (r) => r.timings.duration < 2000,
  });
}

function getProductList() {
  const url = `${BASE_URL}${PRODUCT_LIST_ENDPOINT}?limit=${PRODUCT_LIMIT}`;
  const res = http.get(url, workloadParams);
  checkResponse(res, 'product list');
}

function getProductListByCategory() {
  const categoryId = randomChoice(categoryIds);
  const url = `${BASE_URL}${PRODUCT_LIST_ENDPOINT}?category_id=${categoryId}&limit=${PRODUCT_LIMIT}`;
  const res = http.get(url, workloadParams);
  checkResponse(res, 'product list by category');
}

function getProductDetail() {
  const productId = randomChoice(productIds);
  const url = `${BASE_URL}${PRODUCT_DETAIL_ENDPOINT}/${productId}`;
  const res = http.get(url, workloadParams);
  checkResponse(res, 'product detail');
}

export default function () {
  if (!productIds || productIds.length === 0) {
    throw new Error('productIds is empty. Paste product UUIDs into tests/k6/fixtures/catalog-fixture.js');
  }

  const action = randomInt(1, 100);

  if (action <= 50) {
    getProductList();
  } else if (action <= 80) {
    getProductListByCategory();
  } else {
    getProductDetail();
  }

  sleep(1);
}
