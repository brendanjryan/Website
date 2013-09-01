(function(){

  var changeTime = 5000;
  $strings = $('.infoStrings');

  currentString = -1;
  showNextString = function(){
    ++currentString;
    $strings.eq( currentString % $strings.length)
    .fadeIn(3000)
    .delay(5000)
    .fadeOut(3000, showNextString);
  }

  showNextString();
})();

