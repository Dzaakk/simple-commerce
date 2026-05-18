import { sleep } from 'k6';

import { thinkTime } from './lib/common.js';
import { catalogBrowsingFlow, metricsProbe, setupCommerceData } from './lib/flows.js';

const BASELINE_VUS = Number(__ENV.BASELINE_VUS || 10);
const SPIKE_VUS = Number(__ENV.SPIKE_VUS || 500);

export const options = {
  scenarios: {
    spike_catalog_browse: {
      executor: 'ramping-vus',
      stages: [
        { duration: __ENV.WARMUP || '30s', target: BASELINE_VUS },
        { duration: __ENV.SPIKE_RAMP || '15s', target: SPIKE_VUS },
        { duration: __ENV.SPIKE_HOLD || '1m', target: SPIKE_VUS },
        { duration: __ENV.RECOVERY_RAMP || '15s', target: BASELINE_VUS },
        { duration: __ENV.RECOVERY_HOLD || '1m', target: BASELINE_VUS },
        { duration: __ENV.RAMP_DOWN || '30s', target: 0 },
      ],
      gracefulRampDown: '30s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.10'],
    http_req_duration: ['p(95)<2000', 'p(99)<5000'],
    checks: ['rate>0.90'],
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

  thinkTime(0.1, 0.8);
  sleep(0.05);
}
