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

const DEFAULT_CATALOG_API_VERSION = 'v1';
const SETUP_PRODUCT_LIMIT = 20;
const BROWSE_PRODUCT_LIMIT = 10;

export const CATALOG_API_VERSION = resolveCatalogApiVersion(__ENV);

function resolveCatalogApiVersion(env) {
  return normalizeCatalogApiVersion(env.CATALOG_API_VERSION || env.API_VERSION);
}

function normalizeCatalogApiVersion(version) {
  const normalized = String(version || DEFAULT_CATALOG_API_VERSION).toLowerCase();

  if (normalized === '2' || normalized === 'v2') {
    return 'v2';
  }

  return 'v1';
}

function flowApiVersion(data, apiVersion) {
  return normalizeCatalogApiVersion(
    apiVersion || (data && data.catalogApiVersion) || CATALOG_API_VERSION
  );
}

function catalogUrl(path, apiVersion) {
  return url(`/api/${normalizeCatalogApiVersion(apiVersion)}${path}`);
}

function fetchCatalogJson(path, apiVersion) {
  const res = http.get(catalogUrl(path, apiVersion));

  return {
    res,
    body: parseJson(res),
  };
}

function productListPath(limit, extraParams = {}) {
  const params = {
    limit,
    sort_by: 'newest',
  };

  for (const [key, value] of Object.entries(extraParams)) {
    params[key] = value;
  }

  return `/product${encodeQuery(params)}`;
}

function productItems(body) {
  return hasProductItems(body) ? body.data.items : [];
}

function categoryItems(body) {
  return hasCategories(body) ? body.data : [];
}

function hasProductItems(body) {
  return body && body.data && Array.isArray(body.data.items);
}

function hasCategories(body) {
  return body && Array.isArray(body.data);
}

export function setupCommerceData(apiVersion = CATALOG_API_VERSION) {
  const version = normalizeCatalogApiVersion(apiVersion);

  assertReady();

  const products = fetchCatalogJson(productListPath(SETUP_PRODUCT_LIMIT), version);
  check(products.res, {
    'setup product list status is 200': (r) => r.status === 200,
    'setup product list has success envelope': () => hasSuccessEnvelope(products.body),
  });

  const categories = fetchCatalogJson('/category', version);
  check(categories.res, {
    'setup category list status is 200': (r) => r.status === 200,
    'setup category list has success envelope': () => hasSuccessEnvelope(categories.body),
  });

  return {
    catalogApiVersion: version,
    productIds: productItems(products.body).map((item) => item.id).filter(Boolean),
    categoryIds: categoryItems(categories.body).map((item) => item.ID || item.id).filter(Boolean),
  };
}

export function catalogBrowsingFlow(data, apiVersion) {
  const version = flowApiVersion(data, apiVersion);

  group(`catalog browsing ${version}`, () => {
    browseProductList(version);
    thinkTime();
    browseCategories(version);
  });

  maybeRun(0.6, () => productDetailFlow(data, version));
  maybeRun(0.35, () => categoryFilterFlow(data, version));
}

export function catalogBrowsingFlowV1(data) {
  catalogBrowsingFlow(data, 'v1');
}

export function catalogBrowsingFlowV2(data) {
  catalogBrowsingFlow(data, 'v2');
}

export function productDetailFlow(data, apiVersion) {
  const productId = randomItem(data && data.productIds);
  if (!productId) {
    return;
  }

  group('product detail', () => {
    const { res, body } = fetchCatalogJson(`/product/${productId}`, apiVersion);

    check(res, {
      'product detail status is 200': (r) => r.status === 200,
      'product detail has success envelope': () => hasSuccessEnvelope(body),
      'product detail id matches request': () => body && body.data && body.data.id === productId,
    });
  });

  thinkTime(0.2, 0.8);
}

export function categoryFilterFlow(data, apiVersion) {
  const categoryId = randomItem(data && data.categoryIds);
  if (!categoryId) {
    return;
  }

  group('category filter', () => {
    const { res, body } = fetchCatalogJson(
      productListPath(BROWSE_PRODUCT_LIMIT, { category_id: categoryId }),
      apiVersion
    );

    check(res, {
      'category filter status is 200': (r) => r.status === 200,
      'category filter has success envelope': () => hasSuccessEnvelope(body),
      'category filter data has items array': () => hasProductItems(body),
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

function browseProductList(apiVersion) {
  const { res, body } = fetchCatalogJson(productListPath(BROWSE_PRODUCT_LIMIT), apiVersion);

  check(res, {
    'catalog product list status is 200': (r) => r.status === 200,
    'catalog product list has success envelope': () => hasSuccessEnvelope(body),
    'catalog product list has items array': () => hasProductItems(body),
  });
}

function browseCategories(apiVersion) {
  const { res, body } = fetchCatalogJson('/category', apiVersion);

  check(res, {
    'catalog category list status is 200': (r) => r.status === 200,
    'catalog category list has success envelope': () => hasSuccessEnvelope(body),
    'catalog category list data is array': () => hasCategories(body),
  });
}

function maybeRun(probability, fn) {
  if (Math.random() < probability) {
    fn();
  }
}
