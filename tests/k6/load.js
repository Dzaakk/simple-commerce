import { sleep } from 'k6';

import { thinkTime } from './lib/common.js';
import { catalogBrowsingFlow, metricsProbe, setupCommerceData } from './lib/flows.js';

const TARGET_VUS = Number(__ENV.TARGET_VUS || 100);
const RAMP_UP = __ENV.RAMP_UP || '1m';
const STEADY_STATE = __ENV.STEADY_STATE || '5m';
const RAMP_DOWN = __ENV.RAMP_DOWN || '1m';

export const options = {
  scenarios: {
    steady_catalog_browse: {
      executor: 'ramping-vus',
      stages: [
        { duration: RAMP_UP, target: TARGET_VUS },
        { duration: STEADY_STATE, target: TARGET_VUS },
        { duration: RAMP_DOWN, target: 0 },
      ],
      gracefulRampDown: '30s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<800', 'p(99)<1500'],
    checks: ['rate>0.99'],
  },
};

export function setup() {
  return setupCommerceData();
}

export default function (data) {
  catalogBrowsingFlow(data);

  if (Math.random() < 0.02) {
    metricsProbe();
  }

  thinkTime(0.5, 1.5);
  sleep(0.1);
}
