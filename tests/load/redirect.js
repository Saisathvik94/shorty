
import http from "k6/http";
import { check } from "k6";
import { Counter } from "k6/metrics";

// Track HTTP response status codes
const status302 = new Counter("status_302");
const status404 = new Counter("status_404");
const status429 = new Counter("status_429");
const status500 = new Counter("status_500");
const statusOther = new Counter("status_other");

export const options = {
  vus: 50,
  duration: "30s",
  maxRedirects: 0,
  discardResponseBodies: true,
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:3000";
const SHORT_CODE = __ENV.SHORT_CODE || "A7T1wU";

export default function () {
  const res = http.get(`${BASE_URL}/${SHORT_CODE}`, {
    redirects: 0,
  });

  // Record the actual response status
  switch (res.status) {
    case 302:
      status302.add(1);
      break;

    case 404:
      status404.add(1);
      break;

    case 429:
      status429.add(1);
      break;

    case 500:
      status500.add(1);
      break;

    default:
      statusOther.add(1);
  }

  // Verify the redirect is successful
  check(res, {
    "redirect successful": (r) =>
      r.status === 302 &&
      Boolean(r.headers["Location"]),
  });
}
