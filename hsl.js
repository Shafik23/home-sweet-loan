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
      monthData.monthlyInterest.toFixed(2) +
      '</td><td>' +
      monthData.monthlyPrincipal.toFixed(2) +
      '</td><td>' +
      monthData.remainingBalance.toFixed(2) +
      '</td></tr>';
  });

  table += '</table>';

  document.getElementById('results').innerHTML =
    '<h3> Monthly Payment: ' +
    data.monthlyPayment.toFixed(2) +
    ', Total Payment: ' +
    data.totalPayment.toFixed(2) +
    ', Total Interest: ' +
    data.totalInterest.toFixed(2) +
    '<br><br></h3>' +
    table;
});
