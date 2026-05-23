import http from 'k6/http';
import { check, sleep } from 'k6';
import { productIds, categoryIds } from './fixtures/catalog-fixture.js';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');
const PRODUCT_LIST_ENDPOINT = __ENV.PRODUCT_LIST_ENDPOINT || '/api/v1/product';
const PRODUCT_DETAIL_ENDPOINT = __ENV.PRODUCT_DETAIL_ENDPOINT || '/api/v1/product';
const PRODUCT_LIMIT = Number(__ENV.PRODUCT_LIMIT || 100);

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<2000'],
    checks: ['rate>0.99'],
  },
};

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
  const res = http.get(url);
  checkResponse(res, 'product list');
}

function getProductListByCategory() {
  const categoryId = randomChoice(categoryIds);
  const url = `${BASE_URL}${PRODUCT_LIST_ENDPOINT}?category_id=${categoryId}&limit=${PRODUCT_LIMIT}`;
  const res = http.get(url);
  checkResponse(res, 'product list by category');
}

function getProductDetail() {
  const productId = randomChoice(productIds);
  const url = `${BASE_URL}${PRODUCT_DETAIL_ENDPOINT}/${productId}`;
  const res = http.get(url);
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
