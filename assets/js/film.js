'use strict';

let url = window.location.href;
const regex = /"row":[0-9],|"column":[0-9]/;
window.addEventListener('load', event => {
  let response = fetch(url + '/seats', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
    },
  })
    .then(function (response) {
      return response.text();
    })
    .then(function (data) {
      let seatsRows = JSON.parse(data);
      return seatsRows;
    })
    .then(function initialize(seatsRows) {
      let hasSeatsTaken = false;

      if (seatsRows.Seats != null) {
        hasSeatsTaken = true;
      } else {
        hasSeatsTaken = false;
      }
      initializeSeats(hasSeatsTaken, seatsRows);
    });
});

function initializeSeats(hasSeatsTaken, seatsRows) {
  //Инициализация всех мест с номером ряда и колонки
  let rows = document.querySelectorAll('.seats__row');
  rows.forEach(function (row, index) {
    row.number = index + 1;
    let seats = row.querySelectorAll('.seats__block');
    seats.forEach(function (seat, column) {
      seat.row = row.number;
      seat.column = column + 1;
      if (hasSeatsTaken) {
        seatsRows.Seats.forEach(element => {
          if (seat.row == element.row && seat.column == element.column) {
            seat.classList.toggle('unavailable');
          }
        });
      }
      // seat.price = get from server...
      seat.price = 300;

      if (seat.classList.contains('unavailable')) {
        console.log('unavailable seat');
      } else {
        seat.onclick = handleInterface.bind(seat, seat.row, seat.column);
      }
    });
  });
}

let buyButton = document.querySelector('.purchase_btn');
buyButton.onclick = submitSeats.bind(buyButton);

function handleInterface(row, column) {
  handleCard(this, row, column);
  handleButton();
}
function handleCard(seat, row, column) {
  seat.classList.toggle('chosen');
  let tableCards = document.querySelector('.table__cards');
  let card = document.createElement('card' + row + column);
  card.innerHTML =
    '<div class="card__title">Билет</div><div class="card__seat">Ряд ' +
    row +
    ', место ' +
    column +
    //'...'+ seat.price + '...'
    '</div><div class="card__price">300 Р</div>';

  if (seat.classList.contains('chosen')) {
    console.log('table__cards=', tableCards);
    tableCards.appendChild(card);
    card.classList.toggle('table__card');
    card.classList.toggle('rc' + row + column);
  } else {
    tableCards.removeChild(document.querySelector('.rc' + row + column));
  }

  console.log(column, row);
}

function handleButton() {
  let tablePrice = document.querySelector('.table__price');
  if (checkPrice() == 0) {
    tablePrice.innerHTML = '';
  } else {
    tablePrice.innerHTML = checkPrice() + 'Р';
  }
}

function checkPrice() {
  let price = 0;
  let rows = document.querySelectorAll('.seats__row');
  rows.forEach(function (row, index) {
    row.number = index + 1;
    let seats = row.querySelectorAll('.seats__block');
    seats.forEach(function (seat, column) {
      if (seat.classList.contains('chosen')) {
        price = price + seat.price;
      }
    });
  });
  return price;
}

function submitSeats() {
  // Не инициализировать в цикле с большим количеством выбранных сидений!!!
  let seatsChosen = document.querySelectorAll('.chosen');
  console.log('seatsChosen=', seatsChosen);
  seatsChosen.forEach(function (seat, index) {
    if (index > 10) {
      console.log('Нельзя забронировать более 10 мест!');
    } else {
      let responce = fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8',
        },
        body: JSON.stringify({
          row: seat.row,
          column: seat.column,
          taken: true,
        }),
      });
    }
  });
}
