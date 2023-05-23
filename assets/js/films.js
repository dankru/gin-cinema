let addFilmsButton = document.querySelector('.addOffline-FilmsBtn');
let deleteFilmsButton = document.querySelector('.deleteOffline-FilmsBtn');

let titleInput = document.querySelector('.titleInput');
let textInput = document.querySelector('.textInput');
let timeInput2d = document.querySelector('.timeInput2d');
let priceInput2d = document.querySelector('.priceInput2d');
let timeInput3d = document.querySelector('.timeInput3d');
let priceInput3d = document.querySelector('.priceInput3d');
let imageInput = document.querySelector('.input-img');
let cards = document.querySelectorAll('.card');
let dateInput = document.getElementById('datepicker');

const Title = /\\n +|\\n|^[\s{2,}]+|\d+-\d+-\d+|[\s{2,}]+$| {2,}|\n/g;

addFilmsButton.onclick = AddFilms.bind(
  addFilmsButton,
  titleInput,
  textInput,
  imageInput,
  timeInput2d,
  timeInput3d,
  priceInput2d,
  priceInput3d,
  dateInput,
);
deleteFilmsButton.addEventListener('click', deleteFilms);

async function AddFilms(title, text, image, time2d, time3d, price2d, price3d, date) {
  let validFilm = /[^\._\-\/\*\(\)\=az&]*/;
  let validTimeRegex = /^[0-9]+-[0-9]+-[0-9]+ ?(?:[01]?\d|2[0-3])(?::[0-5]\d){1,2}$/;
  let validPriceRegex = /^ ?[0-9]+$/;

  imagePath = image.value;
  let imageName = imagePath.replace(/^.*[\\\/]/, '');

  price2dValue = price2d.value.replace(/ /g, '');
  let price2dArray = price2dValue.split(',');

  time2dValue = time2d.value.replace(/ /g, '');
  let time2dArray = time2dValue.split(',');
  time2dArray = time2dArray.map(value => {
    return date.value + ' ' + value;
  });

  price3dValue = price3d.value.replace(/ /g, '');
  let price3dArray = price3dValue.split(',');

  time3dValue = time3d.value.replace(/ /g, '');
  let time3dArray = time3dValue.split(',');
  time3dArray = time3dArray.map(value => {
    return date.value + ' ' + value;
  });

  console.log(time2dArray, price2dArray, time3dArray, price3dArray);

  let isValidFilm =
    validFilm.test(title.value) &&
    validFilm.test(text.value) &&
    title.value != '' &&
    text.value != '';

  let isValidTime2d = time2dArray
    .map(value => {
      return validTimeRegex.test(value);
    })
    .every(elem => {
      return elem == true;
    });

  let isValidTime3d = time3dArray
    .map(value => {
      return validTimeRegex.test(value);
    })
    .every(elem => {
      return elem == true;
    });

  let isValidPrice2d = price2dArray
    .map(value => {
      return validPriceRegex.test(value);
    })
    .every(elem => {
      return elem == true;
    });
  let isValidPrice3d = price3dArray
    .map(value => {
      return validPriceRegex.test(value);
    })
    .every(elem => {
      return elem == true;
    });
  let isValidDate = date.value != '';
  let isImagePresent = imageName != '';

  let isInputCorrect =
    isValidFilm &&
    isValidTime2d &&
    isValidTime3d &&
    isValidPrice2d &&
    isValidPrice3d &&
    isValidDate &&
    isImagePresent;
  console.log(
    'film, time2d, time3d, price2d, price3d valid?= ',
    isValidFilm,
    isValidTime2d,
    isValidTime3d,
    isValidPrice2d,
    isValidPrice3d,
  );
  let isLengthCorrect =
    time2dArray.length == price2dArray.length && time3dArray.length == price3dArray.length;
  if (isInputCorrect && isLengthCorrect) {
    if (!imageName) {
      imageName = 'replacement.jpg';
    }

    let addFilmResponse = await fetch('/films', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
      },
      body: JSON.stringify({
        Title: title.value,
        Description: text.value,
        Image: imageName,
      }),
    });

    let selectFilmResponse = await fetch('/films', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
      },
    })
      .then(async function (response) {
        return response;
      })
      .then(async function (data) {
        return data.text();
      })
      .then(async function (FilmId) {
        console.log('response= ', FilmId);
        for (let i = 0; i < time2dArray.length; i++) {
          let datePrice2dResponse = await fetch('/films/datePrice2d', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json;charset=utf-8',
            },
            body: JSON.stringify({
              Date: time2dArray[i],
              Price: Number(price2dArray[i]),
              FilmId: parseInt(FilmId),
            }),
          });
        }
        for (let i = 0; i < time3dArray.length; i++) {
          let datePrice3dResponse = await fetch('/films/datePrice3d', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json;charset=utf-8',
            },
            body: JSON.stringify({
              Date: time3dArray[i],
              Price: Number(price3dArray[i]),
              FilmId: parseInt(FilmId),
            }),
          });
        }
      });
  } else {
    console.log('Некорректный ввод, input, length= ', isInputCorrect, isLengthCorrect);
  }
}
