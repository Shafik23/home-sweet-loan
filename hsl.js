document.addEventListener('htmx:afterOnLoad', function (event) {
  const data = JSON.parse(event.detail.xhr.responseText);
  if (data.error) {
    const resultsDiv = document.getElementById('results');
    resultsDiv.innerHTML = data.error;
    resultsDiv.style.color = 'red';
    resultsDiv.style.fontWeight = 'bold';
    return;
  }

  var table =
    '<table><tr><th>Month</th><th>Interest</th><th>Principal</th><th>Remaining Balance</th></tr>';

  data.schedule.forEach(function (monthData) {
    table +=
      '<tr><td>' +
      monthData.month +
      '</td><td>' +
      renderNumber(monthData.monthlyInterest) +
      '</td><td>' +
      renderNumber(monthData.monthlyPrincipal) +
      '</td><td>' +
      renderNumber(monthData.remainingBalance) +
      '</td></tr>';
  });

  table += '</table>';

  document.getElementById('results').innerHTML =
    '<h3> Monthly Payment: ' +
    renderNumber(data.monthlyPayment) +
    ', Total Payment: ' +
    renderNumber(data.totalPayment) +
    ', Total Interest: ' +
    renderNumber(data.totalInterest) +
    '<br><br></h3>' +
    table;
});

function renderNumber(number) {
  // Javascript is so quirky! This takes care of the the -0.00 case
  if (number < 0 && number > -0.001) {
    number = 0;
  }

  return number.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
}
