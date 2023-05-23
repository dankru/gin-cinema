(()=>{"use strict";$(document).ready((function(){$(".slider-films__slider").slick({slidesToShow:3,responsive:[{breakpoint:1100,settings:{slidesToShow:2}},{breakpoint:767,settings:{slidesToShow:1}}]}),$(".events__slider").slick({slidesToShow:3,responsive:[{breakpoint:1100,settings:{slidesToShow:2}},{breakpoint:767,settings:{slidesToShow:1}}]}),$(".header__burger").click((function(e){$(".header__burger,.header__menu,.header__icons").toggleClass("active"),$("body").toggleClass("lock")}))})),function(e){let s=new Image;s.onload=s.onerror=function(){!function(e){let s=!0===e?"webp":"no-webp";document.documentElement.classList.add(s)}(2==s.height)},s.src="data:image/webp;base64,UklGRjoAAABXRUJQVlA4IC4AAACyAgCdASoCAAIALmk0mk0iIiIiIgBoSygABc6WWgAA/veff/0PP8bA//LwYAAA"}()})();

// let switchTheme = document.querySelector(".switchTheme")
// switchTheme.onclick = function() {
//   let theme = document.getElementById("theme");
//   let amogusLogo = document.querySelector(".amogusLogo");
//   let amogusLogo_webp = document.querySelector(".amogusLogo-webp");
//   let headerIcon = document.querySelector(".header__icon");
//   let headerIcon_webp = document.querySelector(".header__icon-webp");

//   if (theme.getAttribute("href") == "../../../../assets/css/style.min.css") {
//     theme.href = "../../../../assets/css/style.min.light.css";
//     amogusLogo.src= "../../../../assets/img/AmogusLogo-green.jpg";
//     amogusLogo_webp.srcset= "../../../../assets/img/AmogusLogo-green.webp";
//     headerIcon.src= "../../../../assets/img/login-black.png";
//     headerIcon_webp.srcset = "../../../../assets/img/login-black.webp";
//     localStorage.setItem('theme', this.dataset.theme);
//   }
//   else{
//     theme.href = "../../../../assets/css/style.min.css";
//     amogusLogo.src= "../../../../assets/img/AmogusLogo.jpg";
//     amogusLogo_webp.srcset= "../../../../assets/img/AmogusLogo.webp";
//     headerIcon.src= "../../../../assets/img/login.png";
//     headerIcon_webp.srcset = "../../../../assets/img/login.webp";
//   }
// }

