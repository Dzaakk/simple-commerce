import { sleep } from 'k6';

import { resolveLoadProfile } from './config/load-profiles.js';
import { thinkTime } from './lib/common.js';
import { catalogBrowsingFlow, metricsProbe, setupCommerceData } from './lib/flows.js';

const LOAD_PROFILE = resolveLoadProfile(__ENV);

export const options = {
  scenarios: {
    steady_catalog_browse: {
      executor: 'ramping-vus',
      stages: [
        { duration: LOAD_PROFILE.rampUp, target: LOAD_PROFILE.targetVus },
        { duration: LOAD_PROFILE.steadyState, target: LOAD_PROFILE.targetVus },
        { duration: LOAD_PROFILE.rampDown, target: 0 },
      ],
      gracefulRampDown: LOAD_PROFILE.gracefulRampDown,
      tags: {
        test_type: 'load',
        load_profile: LOAD_PROFILE.name,
        target_vus: String(LOAD_PROFILE.targetVus),
      },
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
