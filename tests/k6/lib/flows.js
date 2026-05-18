import http from 'k6/http';
import { check, group } from 'k6';

import {
  assertReady,
  encodeQuery,
  hasSuccessEnvelope,
  parseJson,
  randomItem,
  thinkTime,
  url,
} from './common.js';

export function setupCommerceData() {
  assertReady();

  const productsRes = http.get(url('/api/v1/product?limit=20&sort_by=newest'));
  const productsBody = parseJson(productsRes);

  check(productsRes, {
    'setup product list status is 200': (r) => r.status === 200,
    'setup product list has success envelope': () => hasSuccessEnvelope(productsBody),
  });

  const categoriesRes = http.get(url('/api/v1/category'));
  const categoriesBody = parseJson(categoriesRes);

  check(categoriesRes, {
    'setup category list status is 200': (r) => r.status === 200,
    'setup category list has success envelope': () => hasSuccessEnvelope(categoriesBody),
  });

  const products =
    productsBody && productsBody.data && Array.isArray(productsBody.data.items)
      ? productsBody.data.items
      : [];
  const categories =
    categoriesBody && Array.isArray(categoriesBody.data) ? categoriesBody.data : [];

  return {
    productIds: products.map((item) => item.id).filter(Boolean),
    categoryIds: categories.map((item) => item.ID || item.id).filter(Boolean),
  };
}

export function catalogBrowsingFlow(data) {
  group('catalog browsing', () => {
    const productsRes = http.get(url('/api/v1/product?limit=10&sort_by=newest'));
    const productsBody = parseJson(productsRes);

    check(productsRes, {
      'catalog product list status is 200': (r) => r.status === 200,
      'catalog product list has success envelope': () => hasSuccessEnvelope(productsBody),
      'catalog product list has items array': () =>
        productsBody && productsBody.data && Array.isArray(productsBody.data.items),
    });

    thinkTime();

    const categoriesRes = http.get(url('/api/v1/category'));
    const categoriesBody = parseJson(categoriesRes);

    check(categoriesRes, {
      'catalog category list status is 200': (r) => r.status === 200,
      'catalog category list has success envelope': () => hasSuccessEnvelope(categoriesBody),
      'catalog category list data is array': () =>
        categoriesBody && Array.isArray(categoriesBody.data),
    });
  });

  if (Math.random() < 0.6) {
    productDetailFlow(data);
  }

  if (Math.random() < 0.35) {
    categoryFilterFlow(data);
  }
}

export function productDetailFlow(data) {
  const productId = randomItem(data.productIds);
  if (!productId) {
    return;
  }

  group('product detail', () => {
    const res = http.get(url(`/api/v1/product/${productId}`));
    const body = parseJson(res);

    check(res, {
      'product detail status is 200': (r) => r.status === 200,
      'product detail has success envelope': () => hasSuccessEnvelope(body),
      'product detail id matches request': () => body && body.data && body.data.id === productId,
    });
  });

  thinkTime(0.2, 0.8);
}

export function categoryFilterFlow(data) {
  const categoryId = randomItem(data.categoryIds);
  if (!categoryId) {
    return;
  }

  group('category filter', () => {
    const query = encodeQuery({
      category_id: categoryId,
      limit: 10,
      sort_by: 'newest',
    });
    const res = http.get(url(`/api/v1/product${query}`));
    const body = parseJson(res);

    check(res, {
      'category filter status is 200': (r) => r.status === 200,
      'category filter has success envelope': () => hasSuccessEnvelope(body),
      'category filter data has items array': () =>
        body && body.data && Array.isArray(body.data.items),
    });
  });

  thinkTime(0.2, 0.8);
}

export function metricsProbe() {
  group('metrics probe', () => {
    const res = http.get(url('/metrics'));

    check(res, {
      'metrics probe status is 200': (r) => r.status === 200,
      'metrics probe exposes prometheus text': (r) =>
        r.body.includes('# HELP') && r.body.includes('go_goroutines'),
    });
  });
}
