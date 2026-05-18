import { sleep } from 'k6';

import { thinkTime } from './lib/common.js';
import { catalogBrowsingFlow, metricsProbe, setupCommerceData } from './lib/flows.js';

const MAX_VUS = Number(__ENV.MAX_VUS || 300);
const STEP_1 = Math.max(1, Math.round(MAX_VUS * 0.33));
const STEP_2 = Math.max(STEP_1, Math.round(MAX_VUS * 0.66));

export const options = {
  scenarios: {
    stress_catalog_browse: {
      executor: 'ramping-vus',
      stages: [
        { duration: __ENV.RAMP_UP || '1m', target: STEP_1 },
        { duration: __ENV.STEP_HOLD || '2m', target: STEP_1 },
        { duration: __ENV.RAMP_STEP || '1m', target: STEP_2 },
        { duration: __ENV.STEP_HOLD || '2m', target: STEP_2 },
        { duration: __ENV.RAMP_STEP || '1m', target: MAX_VUS },
        { duration: __ENV.PEAK_HOLD || '3m', target: MAX_VUS },
        { duration: __ENV.RAMP_DOWN || '2m', target: 0 },
      ],
      gracefulRampDown: '45s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<1500', 'p(99)<3000'],
    checks: ['rate>0.95'],
  },
};

export function setup() {
  return setupCommerceData();
}

export default function (data) {
  catalogBrowsingFlow(data);

  if (Math.random() < 0.01) {
    metricsProbe();
  }

  thinkTime(0.2, 1.0);
  sleep(0.05);
}
