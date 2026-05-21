const API = 'http://3.122.83.90:4080/api';
let token = localStorage.getItem('token');
let selectedStars = 0;

function getToken() { return localStorage.getItem('token'); }

function api(path, options = {}) {
  const headers = { 'Content-Type': 'application/json' };
  const t = getToken();
  if (t) headers['Authorization'] = 'Bearer ' + t;
  return fetch(API + path, { ...options, headers }).then(r => r.json());
}

function showMessage(area, text, type) {
  const el = typeof area === 'string' ? document.getElementById(area) : area;
  el.innerHTML = '<div class="message ' + type + '">' + text + '</div>';
  if (type !== 'error') setTimeout(() => el.innerHTML = '', 3000);
}

function toggleAuth() {
  document.getElementById('login-form').classList.toggle('hidden');
  document.getElementById('register-form').classList.toggle('hidden');
  document.getElementById('message-area').innerHTML = '';
}

async function login() {
  const email = document.getElementById('login-email').value;
  const password = document.getElementById('login-password').value;
  const data = await api('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  });
  if (data.token) {
    localStorage.setItem('token', data.token);
    token = data.token;
    initApp();
  } else {
    showMessage('message-area', data.error || 'Giriş hatası', 'error');
  }
}

async function register() {
  const name = document.getElementById('register-name').value;
  const email = document.getElementById('register-email').value;
  const password = document.getElementById('register-password').value;
  const data = await api('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password })
  });
  if (data.token) {
    localStorage.setItem('token', data.token);
    token = data.token;
    initApp();
  } else {
    showMessage('message-area', data.error || 'Kayıt hatası', 'error');
  }
}

const pageIds = ['dashboard', 'calendar', 'profile'];

function logout() {
  localStorage.removeItem('token');
  token = null;
  document.getElementById('navbar').classList.add('hidden');
  document.getElementById('auth-page').classList.remove('hidden');
  pageIds.forEach(p => document.getElementById(p + '-page').classList.add('hidden'));
}

function showPage(page) {
  localStorage.setItem('lastPage', page);
  pageIds.forEach(p => {
    document.getElementById(p + '-page').classList.toggle('hidden', p !== page);
  });
  if (page === 'dashboard') loadTodos();
  if (page === 'calendar') renderCalendar();
  if (page === 'profile') loadProfile();
}

function initApp() {
  document.getElementById('navbar').classList.remove('hidden');
  document.getElementById('auth-page').classList.add('hidden');
  initSubjects();
  showPage('dashboard');
}

function initSubjects() {
  const examSelect = document.getElementById('study-exam');
  const subjectSelect = document.getElementById('study-subject');
  examSelect.value = 'TYT';
  updateSubjects();
  updateTopics();
}

function updateSubjects() {
  const exam = document.getElementById('study-exam').value;
  const subjectSelect = document.getElementById('study-subject');
  subjectSelect.innerHTML = '';
  Object.keys(subjects[exam].topics).forEach(s => {
    const opt = document.createElement('option');
    opt.value = s; opt.textContent = s;
    subjectSelect.appendChild(opt);
  });
  const topicSelect = document.getElementById('study-topic');
  topicSelect.innerHTML = '';
  const first = subjects[exam].topics[subjectSelect.value] || [];
  first.forEach(t => {
    const opt = document.createElement('option');
    opt.value = t; opt.textContent = t;
    topicSelect.appendChild(opt);
  });
}

function updateTopics() {
  const exam = document.getElementById('study-exam').value;
  const subject = document.getElementById('study-subject').value;
  const topicSelect = document.getElementById('study-topic');
  topicSelect.innerHTML = '';
  const topics = subjects[exam].topics[subject] || [];
  topics.forEach(t => {
    const opt = document.createElement('option');
    opt.value = t; opt.textContent = t;
    topicSelect.appendChild(opt);
  });
}

function toggleStudyFields() {
  const type = document.getElementById('study-type').value;
  document.getElementById('stars-field').classList.toggle('hidden', type !== 'goz_gezdir');
  document.getElementById('test-fields').classList.toggle('hidden', type === 'goz_gezdir');
}

function setStars(val) {
  selectedStars = val;
  document.querySelectorAll('.star').forEach(s => {
    s.classList.toggle('active', parseInt(s.dataset.value) <= val);
  });
}

function calculateNet() {
  const correct = parseInt(document.getElementById('test-correct').value) || 0;
  const wrong = parseInt(document.getElementById('test-wrong').value) || 0;
  const net = correct - wrong * 0.25;
  document.getElementById('net-result').textContent = 'Net: ' + net.toFixed(2);
}

async function saveStudy() {
  const body = {
    subject: document.getElementById('study-subject').value,
    topic: document.getElementById('study-topic').value,
    study_type: document.getElementById('study-type').value,
    stars: selectedStars,
    correct: parseInt(document.getElementById('test-correct').value) || 0,
    wrong: parseInt(document.getElementById('test-wrong').value) || 0
  };
  if (body.subject === '') { showMessage('study-message', 'Lütfen ders seçin', 'error'); return; }
  if (body.topic === '') { showMessage('study-message', 'Lütfen konu seçin', 'error'); return; }

  const data = await api('/study', {
    method: 'POST',
    body: JSON.stringify(body)
  });
  if (data.message) {
    showMessage('study-message', '✓ ' + data.message, 'success');
    selectedStars = 0;
    document.querySelectorAll('.star').forEach(s => s.classList.remove('active'));
    document.getElementById('test-correct').value = 0;
    document.getElementById('test-wrong').value = 0;
    document.getElementById('net-result').textContent = 'Net: 0.00';
    loadTodos();
  } else {
    showMessage('study-message', data.error || 'Kayıt hatası', 'error');
  }
}

async function loadTodos() {
  const data = await api('/todos/today');
  const list = document.getElementById('todos-list');
  const empty = document.getElementById('no-todos');
  list.innerHTML = '';
  document.getElementById('today-date').textContent = data.date || '';

  if (!data.todos || data.todos.length === 0) {
    empty.classList.remove('hidden');
    return;
  }
  empty.classList.add('hidden');

  const badgeMap = { goz_gezdir: 'badge-goz', test_coz: 'badge-test', genel_test: 'badge-genel' };
  const labelMap = { goz_gezdir: 'Göz Gezdir', test_coz: 'Test Çöz', genel_test: 'Genel Test' };

  data.todos.forEach(t => {
    const div = document.createElement('div');
    div.className = 'todo-item';
    div.innerHTML = `
      <div class="todo-info">
        <div class="todo-subject">${t.subject} <span class="todo-badge ${badgeMap[t.todo_type] || ''}">${labelMap[t.todo_type] || t.todo_type}</span></div>
        <div class="todo-topic">${t.topic}</div>
      </div>
      <button class="success" style="padding:6px 16px;font-size:0.85rem;" onclick="completeTodo(${t.id})">Tamamla</button>
    `;
    list.appendChild(div);
  });
}

async function completeTodo(id) {
  const data = await api('/todos/' + id + '/complete', { method: 'POST' });
  if (data.message) loadTodos();
}

async function loadProfile() {
  const data = await api('/profile');
  if (data.user) {
    const u = data.user;
    document.getElementById('profile-info').innerHTML = `
      <p><strong>Ad Soyad:</strong> ${u.name}</p>
      <p><strong>Email:</strong> ${u.email}</p>
      <p><strong>Kayıt:</strong> ${formatDate(u.created_at)}</p>
      <p><strong>Silinme:</strong> ${formatDate(u.purge_at)}</p>
      <p><strong>Telegram:</strong> ${(u.telegram_chat_id || 0) > 0 ? 'Bağlı' : 'Bağlı değil'}</p>
    `;
    const hasTelegram = (u.telegram_chat_id || 0) > 0;
    document.getElementById('activation-code').textContent = u.activation_code || (hasTelegram ? 'Bağlı (kod gizli)' : 'Kod oluşturun');
    document.getElementById('disconnect-telegram-btn').style.display = hasTelegram ? 'inline-block' : 'none';
  }
}

async function disconnectTelegram() {
  const data = await api('/auth/disconnect-telegram', { method: 'POST' });
  if (data.message) {
    showMessage('profile-message', '✓ Telegram bağlantısı kesildi', 'success');
    loadProfile();
  } else {
    showMessage('profile-message', data.error || 'Hata', 'error');
  }
}

async function generateActivationCode() {
  const data = await api('/auth/activation-code', { method: 'POST' });
  if (data.activation_code) {
    document.getElementById('activation-code').textContent = data.activation_code;
    showMessage('profile-message', '✓ Aktivasyon kodu oluşturuldu', 'success');
  }
}

let calYear, calMonth;

function renderCalendar(year, month) {
  const now = new Date();
  calYear = year || now.getFullYear();
  calMonth = month || now.getMonth() + 1;

  document.getElementById('cal-title').textContent =
    new Date(calYear, calMonth - 1).toLocaleString('tr-TR', { month: 'long', year: 'numeric' });

  api('/calendar?year=' + calYear + '&month=' + calMonth).then(data => {
    buildCalendarGrid(data.days || []);
  });
}

function buildCalendarGrid(days) {
  const grid = document.getElementById('cal-grid');
  grid.innerHTML = '';

  const dayNames = ['Pzt', 'Sal', 'Çar', 'Per', 'Cum', 'Cmt', 'Paz'];
  dayNames.forEach(n => {
    const d = document.createElement('div');
    d.className = 'cal-day-name';
    d.textContent = n;
    grid.appendChild(d);
  });

  const firstDay = new Date(calYear, calMonth - 1, 1).getDay();
  const startOffset = firstDay === 0 ? 6 : firstDay - 1;
  const daysInMonth = new Date(calYear, calMonth, 0).getDate();
  const daysInPrev = new Date(calYear, calMonth - 1, 0).getDate();
  const todayStr = new Date().toISOString().slice(0, 10);

  const dayMap = {};
  days.forEach(d => { dayMap[d.date] = d; });

  for (let i = 0; i < startOffset; i++) {
    const d = document.createElement('div');
    d.className = 'cal-day other-month';
    d.textContent = daysInPrev - startOffset + 1 + i;
    grid.appendChild(d);
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = calYear + '-' + String(calMonth).padStart(2, '0') + '-' + String(day).padStart(2, '0');
    const el = document.createElement('div');
    el.className = 'cal-day';
    if (dateStr === todayStr) el.classList.add('today');
    const info = dayMap[dateStr];
    if (info) {
      if (info.sessions > 0) el.classList.add('has-data');
      if (info.todos > 0) el.classList.add('has-todo');
      let label = '';
      if (info.sessions > 0) label += info.sessions + ' ders';
      if (info.todos > 0) label += (label ? ' ' : '') + '📋' + info.todos;
      if (label) el.innerHTML = day + '<span class="day-stats">' + label + '</span>';
      else el.textContent = day;
    } else {
      el.textContent = day;
    }
    el.onclick = () => showDayDetail(dateStr, info);
    grid.appendChild(el);
  }

  const totalCells = startOffset + daysInMonth;
  const remaining = (7 - (totalCells % 7)) % 7;
  for (let i = 1; i <= remaining; i++) {
    const d = document.createElement('div');
    d.className = 'cal-day other-month';
    d.textContent = i;
    grid.appendChild(d);
  }

  document.getElementById('cal-day-detail').classList.add('hidden');
  document.getElementById('cal-day-content').innerHTML = '';
}

function showDayDetail(dateStr, dayData) {
  const detail = document.getElementById('cal-day-detail');
  const title = document.getElementById('cal-day-title');
  const content = document.getElementById('cal-day-content');

  title.textContent = new Date(dateStr + 'T12:00:00').toLocaleDateString('tr-TR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });

  content.innerHTML = '<div class="cal-session-list"></div><div style="margin-top:12px;text-align:center;"><em style="color:#888;font-size:0.85rem;">Yükleniyor...</em></div>';
  detail.classList.remove('hidden');

  const list = content.querySelector('.cal-session-list');

  api('/calendar/day?date=' + dateStr).then(data => {
    list.innerHTML = '';
    const labelMap = { goz_gezdir: 'Göz Gezdir', test_coz: 'Test Çöz', genel_test: 'Genel Test' };
    const badgeClass = { goz_gezdir: 'badge-goz', test_coz: 'badge-test', genel_test: 'badge-genel' };

    if (data.sessions && data.sessions.length > 0) {
      const hdr = document.createElement('div');
      hdr.className = 'cal-section-header';
      hdr.textContent = '📚 Çalışma Kayıtları';
      list.appendChild(hdr);
      data.sessions.forEach(s => {
        const div = document.createElement('div');
        div.className = 'cal-session-item';
        let extra = '';
        if (s.study_type === 'goz_gezdir') {
          extra = '⭐'.repeat(s.stars || 0);
        } else {
          extra = s.correct + 'D ' + s.wrong + 'Y = ' + s.net.toFixed(2) + ' net';
        }
        div.innerHTML = '<div class="cal-session-subject">' + s.subject + ' → ' + s.topic + ' <span class="todo-badge ' + (badgeClass[s.study_type] || '') + '">' + (labelMap[s.study_type] || s.study_type) + '</span></div><div class="cal-session-detail">' + extra + '</div>';
        list.appendChild(div);
      });
    }

    api('/calendar/day/todos?date=' + dateStr).then(tdata => {
      if (tdata.todos && tdata.todos.length > 0) {
        const hdr = document.createElement('div');
        hdr.className = 'cal-section-header';
        hdr.textContent = '📋 Yapılacaklar';
        list.appendChild(hdr);
        tdata.todos.forEach(t => {
          const div = document.createElement('div');
          div.className = 'cal-session-item';
          const completed = t.completed;
          let statusHtml = '';
          if (completed) {
            statusHtml = '<span style="color:#28a745;font-weight:600;">✓ Yapıldı</span>';
          } else {
            statusHtml = '<button class="success" style="padding:4px 12px;font-size:0.8rem;" onclick="completeTodoFromCal(' + t.id + ',\'' + dateStr + '\')">Tamamla</button>';
          }
          div.innerHTML = '<div class="cal-session-subject">' + t.subject + ' → ' + t.topic + ' <span class="todo-badge ' + (badgeClass[t.todo_type] || '') + '">' + (labelMap[t.todo_type] || t.todo_type) + '</span></div><div class="cal-session-detail" style="margin-top:6px;">' + statusHtml + '</div>';
          list.appendChild(div);
        });
      }

      if ((!data.sessions || data.sessions.length === 0) && (!tdata.todos || tdata.todos.length === 0)) {
        list.innerHTML = '<p style="color:#888;text-align:center;padding:20px;">Bu günde kayıt bulunmuyor.</p>';
      }
      const em = content.querySelector('em');
      if (em) em.remove();
    });
  });
}

function completeTodoFromCal(id, dateStr) {
  api('/todos/' + id + '/complete', { method: 'POST' }).then(data => {
    if (data.message) showDayDetail(dateStr, null);
  });
}

function calendarPrevMonth() {
  if (calMonth === 1) { calYear--; calMonth = 12; }
  else calMonth--;
  renderCalendar(calYear, calMonth);
}

function calendarNextMonth() {
  if (calMonth === 12) { calYear++; calMonth = 1; }
  else calMonth++;
  renderCalendar(calYear, calMonth);
}

function formatDate(d) {
  if (!d) return '-';
  const date = new Date(d);
  return date.toLocaleDateString('tr-TR');
}

document.addEventListener('DOMContentLoaded', () => {
  if (token) {
    initApp();
  } else {
    document.getElementById('auth-page').classList.remove('hidden');
  }
});
