export const LOAD_PROFILES = {
  100: {
    targetVus: 100,
    rampUp: '1m',
    steadyState: '5m',
    rampDown: '1m',
    gracefulRampDown: '30s',
  },
  300: {
    targetVus: 300,
    rampUp: '2m',
    steadyState: '5m',
    rampDown: '2m',
    gracefulRampDown: '45s',
  },
  500: {
    targetVus: 500,
    rampUp: '3m',
    steadyState: '7m',
    rampDown: '3m',
    gracefulRampDown: '1m',
  },
  1000: {
    targetVus: 1000,
    rampUp: '5m',
    steadyState: '10m',
    rampDown: '5m',
    gracefulRampDown: '2m',
  },
};

export function resolveLoadProfile(env) {
  const requestedProfile = env.LOAD_PROFILE || '100';
  const baseProfile = LOAD_PROFILES[requestedProfile] || LOAD_PROFILES[100];

  return {
    name: LOAD_PROFILES[requestedProfile] ? requestedProfile : '100',
    targetVus: Number(env.TARGET_VUS || baseProfile.targetVus),
    rampUp: env.RAMP_UP || baseProfile.rampUp,
    steadyState: env.STEADY_STATE || baseProfile.steadyState,
    rampDown: env.RAMP_DOWN || baseProfile.rampDown,
    gracefulRampDown: env.GRACEFUL_RAMP_DOWN || baseProfile.gracefulRampDown,
  };
}
