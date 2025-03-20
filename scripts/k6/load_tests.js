import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

export let options = {
    stages: [
        { duration: '1m', target: 1000 }, // Рамп-ап: до 1000 виртуальных пользователей за 1 минуту
        { duration: '3m', target: 1000 }, // Стабильная нагрузка: 1000 VU в течение 3 минут
        { duration: '1m', target: 0 },    // Рамп-даун: снижение до 0 VU за 1 минуту
    ],
    thresholds: {
        // 95-й перцентиль времени ответа должен быть меньше 50 мс
        'http_req_duration': ['p(95)<50'],
        // Процент неуспешных запросов должен быть меньше 0.01%
        'http_req_failed': ['rate<0.0001'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const items = [
    't-shirt', 'cup', 'book', 'pen',
    'powerbank', 'hoody', 'umbrella',
    'socks', 'wallet', 'pink-hoody'
];

export default function () {
    let username = `user_${__VU}_${__ITER}`;
    let password = 'testpassword';

    // 1. Регистрация
    let res = http.post(`${BASE_URL}/api/auth/register`, JSON.stringify({
        username: username,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });
    check(res, {
        'registration status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    });

    // 2. Логин — получаем JWT-токен
    res = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
        username: username,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });
    check(res, {
        'login status is 200': (r) => r.status === 200,
    });
    let token = res.json('token');
    let authHeaders = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
    };

    // 3. GET запрос: /api/info
    res = http.get(`${BASE_URL}/api/info`, authHeaders);
    check(res, { 'GET /info status is 200': (r) => r.status === 200 });



    // 4. POST запрос: /api/send-coin (передача монет)
    let receiver_id = Math.floor(Math.random() * 20) + 1;
    res = http.post(`${BASE_URL}/api/send-coin`, JSON.stringify({
        receiver_id: receiver_id,
        amount: 10,
    }), authHeaders);
    check(res, { 'POST /send-coin status is 200': (r) => r.status === 200 });

    // 5. GET запрос: /api/merch
    res = http.get(`${BASE_URL}/api/merch`, authHeaders);
    check(res, { 'GET /merch status is 200': (r) => r.status === 200 });

    // 6. POST запрос: /api/merch/buy/:item
    const randomItem = items[Math.floor(Math.random() * items.length)];

    res = http.post(`${BASE_URL}/api/merch/buy/${randomItem}`, null, authHeaders);
    check(res, { 'POST /merch/buy/:item status is 200': (r) => r.status === 200 });

    // Пауза между итерациями
    sleep(3);
}
