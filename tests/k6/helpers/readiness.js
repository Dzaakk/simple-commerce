import http from 'k6/http';
import { sleep } from 'k6';

const DEFAULT_MAX_ATTEMPTS = 30;
const DEFAULT_RETRY_INTERVAL_SECONDS = 1;
const DEFAULT_REQUEST_TIMEOUT = '2s';

function positiveInteger(value, fallback) {
  const parsed = Number(value);
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback;
}

function positiveNumber(value, fallback) {
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
}

export function waitForReadiness(baseUrl) {
  const maxAttempts = positiveInteger(
    __ENV.READINESS_MAX_ATTEMPTS,
    DEFAULT_MAX_ATTEMPTS
  );
  const retryIntervalSeconds = positiveNumber(
    __ENV.READINESS_RETRY_INTERVAL_SECONDS,
    DEFAULT_RETRY_INTERVAL_SECONDS
  );
  const requestTimeout = __ENV.READINESS_REQUEST_TIMEOUT || DEFAULT_REQUEST_TIMEOUT;
  const readinessUrl = `${baseUrl.replace(/\/+$/, '')}/readyz`;

  let lastStatus = 0;
  let lastError = '';

  for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
    const response = http.get(readinessUrl, {
      timeout: requestTimeout,
      tags: {
        name: 'readiness_gate',
        traffic_type: 'readiness',
      },
    });

    lastStatus = response.status;
    lastError = response.error || '';

    if (response.status === 200) {
      console.info(`API ready after ${attempt} attempt(s): ${readinessUrl}`);
      return;
    }

    console.warn(
      `API not ready (${attempt}/${maxAttempts}): status=${lastStatus} error=${lastError || 'none'}`
    );

    if (attempt < maxAttempts) {
      sleep(retryIntervalSeconds);
    }
  }

  throw new Error(
    `API readiness gate failed after ${maxAttempts} attempts: url=${readinessUrl} status=${lastStatus} error=${lastError || 'none'}`
  );
}
