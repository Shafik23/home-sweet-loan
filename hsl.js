document.addEventListener('DOMContentLoaded', function () {
  document.getElementById('loan-input').addEventListener('htmx:afterOnLoad', function (event) {
    const data = JSON.parse(event.detail.xhr.responseText);
    if (data.error) {
      const resultsDiv = document.getElementById('results');
      resultsDiv.innerHTML = data.error;
      resultsDiv.style.color = 'red';
      resultsDiv.style.fontWeight = 'bold';
      return;
    }

    var table = '<table><tr><th>Month</th><th>Interest</th><th>Principal</th><th>Remaining Balance</th></tr>';

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

  const loanHistoryDropdown = document.getElementById('loanHistoryDropdown');

  if (loanHistoryDropdown) {
    loanHistoryDropdown.addEventListener('htmx:afterOnLoad', function (event) {
      populateLoanHistoryDropdown(JSON.parse(event.detail.xhr.responseText));
    });
  }
});

function populateLoanHistoryDropdown(loans) {
  const dropdown = document.getElementById('loanHistoryDropdown');
  let optionsHtml = '<option value="">-- Select a Loan --</option>';

  loans.forEach(function (loan) {
    optionsHtml += `<option value="${loan.id}">Principal: ${renderNumber(loan.principal)}, Interest Rate: ${
      loan.interestRate
    }%, Term: ${loan.loanTerm} years</option>`;
  });

  dropdown.innerHTML = optionsHtml;
}
function renderNumber(number) {
  // Javascript is so quirky! This takes care of the the -0.00 case
  if (number < 0 && number > -0.001) {
    number = 0;
  }

  return number.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
}
