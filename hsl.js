document.addEventListener('DOMContentLoaded', function () {
  fetchCurrentMarketRate();

  document.getElementById('loan-input').addEventListener('htmx:afterOnLoad', function (event) {
    const data = JSON.parse(event.detail.xhr.responseText);

    if (data.error) {
      const resultsDiv = document.getElementById('results');
      resultsDiv.innerHTML = data.error;
      resultsDiv.style.color = 'red';
      resultsDiv.style.fontWeight = 'bold';
      return;
    }

    const summaryTable =
      '<table>' +
      '<tr><th>Number of Payments</th><th>Monthly Payment</th><th>Total Amount Paid</th><th>Total Interest</th></tr>' +
      '<tr>' +
      '<td>' +
      data.schedule.length +
      '</td>' +
      '<td>' +
      renderNumber(data.monthlyPayment) +
      '</td>' +
      '<td>' +
      renderNumber(data.totalPayment) +
      '</td>' +
      '<td>' +
      renderNumber(data.totalInterest) +
      '</td>' +
      '</tr>' +
      '</table><br>';

    let scheduleTable =
      '<table>' + '<tr><th>Month</th><th>Interest</th><th>Principal</th><th>Remaining Balance</th></tr>';

    data.schedule.forEach(function (monthData) {
      scheduleTable +=
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

    scheduleTable += '</table>';

    document.getElementById('results').innerHTML = summaryTable + scheduleTable;
  });

  const loanHistoryDropdown = document.getElementById('loanHistoryDropdown');

  if (loanHistoryDropdown) {
    // Listen for API data coming back from DB
    loanHistoryDropdown.addEventListener('htmx:afterOnLoad', function (event) {
      populateLoanHistoryDropdown(JSON.parse(event.detail.xhr.responseText));
    });

    // Listen for changes on the loan history dropdown
    loanHistoryDropdown.addEventListener('change', function (event) {
      const selectedText = event.target.selectedOptions[0].text;

      // Extract values from the option's text
      const principalMatch = selectedText.match(/Principal: \$([\d,]+(\.\d{2})?)/);
      const interestRateMatch = selectedText.match(/Interest Rate: (\d+(\.\d{1,2})?)%/);
      const loanTermMatch = selectedText.match(/Term: (\d+) years/);

      // Update input fields if matches were found
      if (principalMatch) {
        const principalValue = principalMatch[1].replace(/,/g, ''); // remove commas
        document.querySelector('[name="principal"]').value = principalValue;
      }
      if (interestRateMatch) {
        document.querySelector('[name="interest_rate"]').value = interestRateMatch[1];
      }
      if (loanTermMatch) {
        document.querySelector('[name="loan_term_years"]').value = loanTermMatch[1];
      }

      // Clear the results elements
      document.getElementById('results').innerHTML = '';
    });
  }
});

// Function to fetch and display the current interest rate
function fetchCurrentMarketRate() {
  fetch('/currentMarketRate')
    .then(response => response.json())
    .then(data => {
      document.getElementById('currentRate').textContent = `Current US Market Rate: ${data.rate}%`;
    })
    .catch(error => {
      console.error('Error fetching current rate:', error);
    });
}

function populateLoanHistoryDropdown(loans) {
  const dropdown = document.getElementById('loanHistoryDropdown');
  let optionsHtml = '';

  if (!loans || loans.length === 0) {
    optionsHtml += '<option value="" disabled>No loans found</option>';
  } else {
    loans.forEach(function (loan) {
      optionsHtml += `<option value="${loan.id}">Principal: ${renderNumber(loan.principal)}, Interest Rate: ${
        loan.interestRate
      }%, Term: ${loan.loanTerm} years</option>`;
    });
  }

  dropdown.innerHTML = optionsHtml;
}

function renderNumber(number) {
  // Javascript is so quirky! This takes care of the the -0.00 case
  if (number < 0 && number > -0.001) {
    number = 0;
  }

  return number.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
}
