import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

const errorRate = new Rate('errors');

const listDuration = new Trend('list_duration', true);
const createDuration = new Trend('create_duration', true);
const getDuration = new Trend('get_duration', true);
const deleteDuration = new Trend('delete_duration', true);

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const headers = { 'Content-Type': 'application/json' };

export const options = {
  scenarios: {
    // Scenario 1: Hammer list/get — fast, CPU-bound, triggers HPA
    read_heavy: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '15s', target: 50 },   // ramp up
        { duration: '2m',  target: 100 },  // sustained high load
        { duration: '15s', target: 0 },    // ramp down
      ],
      exec: 'readHeavy',
    },
    // Scenario 2: Light CRUD — proves full lifecycle works under load
    crud_lifecycle: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '15s', target: 5 },
        { duration: '2m',  target: 10 },
        { duration: '15s', target: 0 },
      ],
      exec: 'crudLifecycle',
    },
  },
  thresholds: {
    'http_req_failed':            ['rate<0.01'],
    'list_duration':              ['p(95)<500'],
    'get_duration':               ['p(95)<500'],
    'errors':                     ['rate<0.01'],
  },
};

// --------------- Scenario 1: Read Heavy ---------------
// Fires list + get as fast as possible to spike CPU and trigger HPA
export function readHeavy() {
  // List all databases
  const listRes = http.get(`${BASE_URL}/databases`, { headers });
  listDuration.add(listRes.timings.duration);
  check(listRes, { 'list: status 200': (r) => r.status === 200 }) || errorRate.add(1);

  // Get a known database (if any exist) or hit the 404 path (still CPU work)
  const getRes = http.get(`${BASE_URL}/databases/loadtest-seed`, { headers });
  getDuration.add(getRes.timings.duration);
  check(getRes, {
    'get: status 200 or 404': (r) => r.status === 200 || r.status === 404,
  }) || errorRate.add(1);

  // Health check (lightweight, adds request volume)
  http.get(`${BASE_URL}/health`);
}

// --------------- Scenario 2: CRUD Lifecycle ---------------
// Full create → get → delete cycle at low concurrency
export function crudLifecycle() {
  const dbName = `loadtest-${Math.random().toString(36).substring(2, 8)}`;

  // Create
  const createRes = http.post(
    `${BASE_URL}/databases`,
    JSON.stringify({ name: dbName, instances: 1, storage: '1Gi' }),
    { headers }
  );
  createDuration.add(createRes.timings.duration);
  const created = check(createRes, { 'create: status 201': (r) => r.status === 201 });
  if (!created) { errorRate.add(1); return; }

  // Get
  const getRes = http.get(`${BASE_URL}/databases/${dbName}`, { headers });
  getDuration.add(getRes.timings.duration);
  check(getRes, { 'get: status 200': (r) => r.status === 200 }) || errorRate.add(1);

  // Delete
  const delRes = http.del(`${BASE_URL}/databases/${dbName}`, null, { headers });
  deleteDuration.add(delRes.timings.duration);
  check(delRes, { 'delete: status 204': (r) => r.status === 204 }) || errorRate.add(1);

  sleep(1); // small pause between CRUD cycles
}

export function handleSummary(data) {
  return {
    stdout: textSummary(data, { indent: '  ', enableColors: true }),
    'load-test-result.json': JSON.stringify(data, null, 2),
  };
}