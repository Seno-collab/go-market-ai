import http from "k6/http";
import { check, sleep } from "k6";
import { Trend, Rate } from "k6/metrics";

const ttfb = new Trend("ttfb");            // thời gian phản hồi (ms)
const errors = new Rate("errors");         // tỉ lệ lỗi

export const options = {
  stages: [
    { duration: "1m", target: 20 }, // tăng dần lên 20 user ảo
    { duration: "3m", target: 20 }, // giữ tải
    { duration: "1m", target: 0 },  // giảm về 0
  ],
  thresholds: {
    errors: ["rate<0.01"],                 // lỗi < 1%
    http_req_duration: ["p(95)<800"],      // p95 < 800ms
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const TOKEN = __ENV.TOKEN; // nếu cần auth

export default function () {
  const params = {
    headers: {
      "Content-Type": "application/json",
      ...(TOKEN ? { Authorization: `Bearer ${TOKEN}` } : {}),
    },
    tags: { name: "GET_/api/health" },
  };

  const res = http.get(`${BASE}/api/health`, params);

  const ok = check(res, {
    "status 200": (r) => r.status === 200,
  });

  errors.add(!ok);
  ttfb.add(res.timings.waiting);

  sleep(1); // “think time” giống người thật
}
