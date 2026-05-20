import { sleep } from 'k6';

import { resolveLoadProfile } from './config/load-profiles.js';
import { thinkTime } from './lib/common.js';
import {
  CATALOG_API_VERSION,
  catalogBrowsingFlow,
  metricsProbe,
  setupCommerceData,
} from './lib/flows.js';

const LOAD_PROFILE = resolveLoadProfile(__ENV);

export const options = {
  scenarios: {
    steady_catalog_browse: catalogBrowseScenario(LOAD_PROFILE),
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<800', 'p(99)<1500'],
    checks: ['rate>0.99'],
  },
};

export function setup() {
  return setupCommerceData(CATALOG_API_VERSION);
}

export default function (data) {
  catalogBrowsingFlow(data);
  maybeProbeMetrics();
  pauseBetweenIterations();
}

function catalogBrowseScenario(profile) {
  return {
    executor: 'ramping-vus',
    stages: rampingStages(profile),
    gracefulRampDown: profile.gracefulRampDown,
    tags: scenarioTags(profile),
  };
}

function rampingStages(profile) {
  return [
    { duration: profile.rampUp, target: profile.targetVus },
    { duration: profile.steadyState, target: profile.targetVus },
    { duration: profile.rampDown, target: 0 },
  ];
}

function scenarioTags(profile) {
  return {
    test_type: 'load',
    load_profile: profile.name,
    target_vus: String(profile.targetVus),
    catalog_api_version: CATALOG_API_VERSION,
  };
}

function maybeProbeMetrics() {
  if (Math.random() < 0.02) {
    metricsProbe();
  }
}

function pauseBetweenIterations() {
  thinkTime(0.5, 1.5);
  sleep(0.1);
}
