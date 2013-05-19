$(document).ready(function(){


  //make icons appear when "connect" button is clicked
  $('.connect').click(function(){

    //reset all animations
    $('.icons').removeClass('animated fadeInDown');
    $('.icons').removeClass('animated fadeOutUp');

   if( $('.icons').hasClass('not-visible') ){ //if the icons hidden class is not active

    console.log("icons hidden -- now fading down")
    $('.icons').removeClass('not-visible');
    $('.icons').addClass('animated fadeInDown');
   } else {
    $('.icons').addClass('animated fadeOutUp');
    //a timeout is needed to allowthe animation to finish before the
    //hidden attribute is applied again
    setTimeout(function(){
      $('.icons').addClass('not-visible');
    }, 400);
    }
  });
});