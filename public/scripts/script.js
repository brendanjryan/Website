$(document).ready(function(){

  $('.connect').click(function(){

    //hard animations reset
    $('.icons').removeClass('animated fadeInDown');
    $('.icons').removeClass('animated fadeOutUp');

   if( $('.icons').hasClass('not-visible') ){
    //fade icons in
    $('.icons').removeClass('not-visible');
    $('.icons').addClass('animated fadeInDown');
   } else {
    $('.icons').addClass('animated fadeOutUp');
    //wait for animation to finish until class is applied
    setTimeout(function(){
      $('.icons').addClass('not-visible');
    }, 400);
    }
  });


});