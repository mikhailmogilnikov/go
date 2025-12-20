/**
 * Google Apps Script для интеграции с Expense Tracker API
 */

// === КОНФИГУРАЦИЯ ===
const API_BASE_URL = 'https://invigoratedly-diaphanometric-kylie.ngrok-free.dev/api';

// Общие заголовки для всех запросов (включая обход страницы ngrok)
function getHeaders(token) {
  const headers = {
    'ngrok-skip-browser-warning': 'true',
    'User-Agent': 'GoogleAppsScript'
  };
  if (token) {
    headers['Authorization'] = 'Bearer ' + token;
  }
  return headers;
}

// === МЕНЮ ===
function onOpen() {
  const ui = SpreadsheetApp.getUi();
  ui.createMenu('💰 Expense Tracker')
    .addItem('🔐 Войти / Регистрация', 'showAuthDialog')
    .addSeparator()
    .addItem('➕ Добавить транзакцию', 'addTransactionFromSheet')
    .addItem('📊 Загрузить транзакции', 'loadTransactions')
    .addSeparator()
    .addItem('💵 Установить бюджет', 'setBudgetFromSheet')
    .addItem('📋 Загрузить бюджеты', 'loadBudgets')
    .addSeparator()
    .addItem('📈 Получить отчёт', 'loadReport')
    .addToUi();
}

// === АВТОРИЗАЦИЯ ===
function showAuthDialog() {
  const html = HtmlService.createHtmlOutput(`
    <style>
      body { font-family: Arial, sans-serif; padding: 20px; }
      input { width: 100%; padding: 8px; margin: 5px 0; box-sizing: border-box; }
      button { width: 100%; padding: 10px; margin: 5px 0; cursor: pointer; }
      .primary { background: #4285f4; color: white; border: none; }
      .secondary { background: #f1f1f1; border: 1px solid #ddd; }
    </style>
    <h3>Авторизация</h3>
    <input type="email" id="email" placeholder="Email">
    <input type="password" id="password" placeholder="Пароль (мин. 6 символов)">
    <button class="primary" onclick="login()">Войти</button>
    <button class="secondary" onclick="register()">Зарегистрироваться</button>
    <p id="status"></p>
    <script>
      function login() {
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        google.script.run
          .withSuccessHandler(r => {
            document.getElementById('status').textContent = r;
            if (r.includes('Успешно')) google.script.host.close();
          })
          .withFailureHandler(e => document.getElementById('status').textContent = e)
          .doLogin(email, password);
      }
      function register() {
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        google.script.run
          .withSuccessHandler(r => {
            document.getElementById('status').textContent = r;
            if (r.includes('Успешно')) google.script.host.close();
          })
          .withFailureHandler(e => document.getElementById('status').textContent = e)
          .doRegister(email, password);
      }
    </script>
  `)
  .setWidth(300)
  .setHeight(250);
  SpreadsheetApp.getUi().showModalDialog(html, 'Вход');
}

function doLogin(email, password) {
  const response = UrlFetchApp.fetch(API_BASE_URL + '/auth/login', {
    method: 'post',
    contentType: 'application/json',
    headers: getHeaders(),
    payload: JSON.stringify({ email, password }),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() === 200) {
    const data = JSON.parse(response.getContentText());
    PropertiesService.getUserProperties().setProperty('AUTH_TOKEN', data.token);
    return 'Успешно! Токен сохранён.';
  }
  return 'Ошибка: ' + response.getContentText();
}

function doRegister(email, password) {
  const response = UrlFetchApp.fetch(API_BASE_URL + '/auth/register', {
    method: 'post',
    contentType: 'application/json',
    headers: getHeaders(),
    payload: JSON.stringify({ email, password }),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() === 201) {
    const data = JSON.parse(response.getContentText());
    PropertiesService.getUserProperties().setProperty('AUTH_TOKEN', data.token);
    return 'Успешно! Вы зарегистрированы.';
  }
  return 'Ошибка: ' + response.getContentText();
}

function getToken() {
  return PropertiesService.getUserProperties().getProperty('AUTH_TOKEN');
}

// === ТРАНЗАКЦИИ ===
function addTransactionFromSheet() {
  const sheet = SpreadsheetApp.getActiveSheet();
  const row = sheet.getActiveRange().getRow();
  
  const amount = sheet.getRange(row, 1).getValue();
  const category = sheet.getRange(row, 2).getValue();
  const description = sheet.getRange(row, 3).getValue() || '';
  const date = sheet.getRange(row, 4).getValue();
  
  if (!amount || !category) {
    SpreadsheetApp.getUi().alert('Заполните сумму (A) и категорию (B)');
    return;
  }
  
  const token = getToken();
  if (!token) {
    SpreadsheetApp.getUi().alert('Сначала войдите в систему');
    return;
  }
  
  const payload = {
    amount: Number(amount),
    category: String(category),
    description: String(description)
  };
  
  if (date) {
    payload.date = Utilities.formatDate(new Date(date), 'GMT', 'yyyy-MM-dd');
  }
  
  const response = UrlFetchApp.fetch(API_BASE_URL + '/transactions', {
    method: 'post',
    contentType: 'application/json',
    headers: getHeaders(token),
    payload: JSON.stringify(payload),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() === 201) {
    const data = JSON.parse(response.getContentText());
    let msg = 'Транзакция добавлена! ID: ' + data.id;
    if (data.budget_warning) {
      msg += '\n⚠️ ' + data.budget_warning;
    }
    SpreadsheetApp.getUi().alert(msg);
  } else {
    SpreadsheetApp.getUi().alert('Ошибка: ' + response.getContentText());
  }
}

function loadTransactions() {
  const token = getToken();
  if (!token) {
    SpreadsheetApp.getUi().alert('Сначала войдите в систему');
    return;
  }
  
  const response = UrlFetchApp.fetch(API_BASE_URL + '/transactions', {
    headers: getHeaders(token),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() !== 200) {
    SpreadsheetApp.getUi().alert('Ошибка: ' + response.getContentText());
    return;
  }
  
  const transactions = JSON.parse(response.getContentText());
  
  let sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName('Транзакции');
  if (!sheet) {
    sheet = SpreadsheetApp.getActiveSpreadsheet().insertSheet('Транзакции');
  }
  sheet.clear();
  
  sheet.getRange(1, 1, 1, 5).setValues([['ID', 'Сумма', 'Категория', 'Описание', 'Дата']]);
  sheet.getRange(1, 1, 1, 5).setFontWeight('bold');
  
  if (transactions.length > 0) {
    const data = transactions.map(tx => [tx.id, tx.amount, tx.category, tx.description, tx.date]);
    sheet.getRange(2, 1, data.length, 5).setValues(data);
  }
  
  SpreadsheetApp.getUi().alert('Загружено ' + transactions.length + ' транзакций');
}

// === БЮДЖЕТЫ ===
function setBudgetFromSheet() {
  const sheet = SpreadsheetApp.getActiveSheet();
  const row = sheet.getActiveRange().getRow();
  
  const category = sheet.getRange(row, 1).getValue();
  const limitAmount = sheet.getRange(row, 2).getValue();
  const period = sheet.getRange(row, 3).getValue() || 'monthly';
  
  if (!category || !limitAmount) {
    SpreadsheetApp.getUi().alert('Заполните категорию (A) и лимит (B)');
    return;
  }
  
  const token = getToken();
  if (!token) {
    SpreadsheetApp.getUi().alert('Сначала войдите в систему');
    return;
  }
  
  const response = UrlFetchApp.fetch(API_BASE_URL + '/budgets', {
    method: 'post',
    contentType: 'application/json',
    headers: getHeaders(token),
    payload: JSON.stringify({
      category: String(category),
      limit_amount: Number(limitAmount),
      period: String(period)
    }),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() === 201) {
    SpreadsheetApp.getUi().alert('Бюджет установлен!');
  } else {
    SpreadsheetApp.getUi().alert('Ошибка: ' + response.getContentText());
  }
}

function loadBudgets() {
  const token = getToken();
  if (!token) {
    SpreadsheetApp.getUi().alert('Сначала войдите в систему');
    return;
  }
  
  const response = UrlFetchApp.fetch(API_BASE_URL + '/budgets', {
    headers: getHeaders(token),
    muteHttpExceptions: true
  });
  
  if (response.getResponseCode() !== 200) {
    SpreadsheetApp.getUi().alert('Ошибка: ' + response.getContentText());
    return;
  }
  
  const budgets = JSON.parse(response.getContentText());
  
  let sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName('Бюджеты');
  if (!sheet) {
    sheet = SpreadsheetApp.getActiveSpreadsheet().insertSheet('Бюджеты');
  }
  sheet.clear();
  
  sheet.getRange(1, 1, 1, 4).setValues([['ID', 'Категория', 'Лимит', 'Период']]);
  sheet.getRange(1, 1, 1, 4).setFontWeight('bold');
  
  if (budgets.length > 0) {
    const data = budgets.map(b => [b.id, b.category, b.limit_amount, b.period]);
    sheet.getRange(2, 1, data.length, 4).setValues(data);
  }
  
  SpreadsheetApp.getUi().alert('Загружено ' + budgets.length + ' бюджетов');
}

// === ОТЧЁТЫ ===
function loadReport() {
  const ui = SpreadsheetApp.getUi();
  
  const fromResult = ui.prompt('Введите начальную дату (YYYY-MM-DD):');
  if (fromResult.getSelectedButton() !== ui.Button.OK) return;
  
  const toResult = ui.prompt('Введите конечную дату (YYYY-MM-DD):');
  if (toResult.getSelectedButton() !== ui.Button.OK) return;
  
  const from = fromResult.getResponseText();
  const to = toResult.getResponseText();
  
  const token = getToken();
  if (!token) {
    ui.alert('Сначала войдите в систему');
    return;
  }
  
  const response = UrlFetchApp.fetch(
    API_BASE_URL + '/reports?from=' + from + '&to=' + to,
    {
      headers: getHeaders(token),
      muteHttpExceptions: true
    }
  );
  
  if (response.getResponseCode() !== 200) {
    ui.alert('Ошибка: ' + response.getContentText());
    return;
  }
  
  const report = JSON.parse(response.getContentText());
  
  let sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName('Отчёт');
  if (!sheet) {
    sheet = SpreadsheetApp.getActiveSpreadsheet().insertSheet('Отчёт');
  }
  sheet.clear();
  
  sheet.getRange(1, 1).setValue('Отчёт по расходам: ' + from + ' - ' + to);
  sheet.getRange(1, 1).setFontWeight('bold').setFontSize(14);
  
  sheet.getRange(2, 1).setValue('Всего потрачено: ' + report.total_expenses);
  
  sheet.getRange(4, 1, 1, 4).setValues([['Категория', 'Потрачено', 'Лимит', '% бюджета']]);
  sheet.getRange(4, 1, 1, 4).setFontWeight('bold');
  
  if (report.categories && report.categories.length > 0) {
    const data = report.categories.map(c => [
      c.category,
      c.total,
      c.budget_limit || '-',
      c.budget_percentage ? c.budget_percentage.toFixed(1) + '%' : '-'
    ]);
    sheet.getRange(5, 1, data.length, 4).setValues(data);
  }
  
  ui.alert('Отчёт загружен!');
}

// === ИНИЦИАЛИЗАЦИЯ ===
function init() {
  onOpen();
  SpreadsheetApp.getUi().alert(
    'Expense Tracker подключён!\n\n' +
    'Используйте меню "💰 Expense Tracker" для работы с системой.'
  );
}
