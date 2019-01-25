function countdate(limitdatestr){

  var now = document.getElementById('now');
  var limitdate = document.getElementById('limitdate');
  var timelimit = document.getElementById('timelimit');

  var limitDate = new Date(limitdatestr);
  limitdate.textContent = limitDate.toLocaleString();
  timer = setInterval(function(){
    var nowDate = new Date();
    now.textContent = nowDate.toLocaleString();
    var rTime = (limitDate - nowDate) / 1000;
    var addZero = function(n) {
      return ('0' + n).slice(-2);
    }
    var gDate = function(rTime) {
      if(rTime > 0.0){
        var rDay = Math.floor(rTime / (60 * 60 * 24));
        var rHour = Math.floor(rTime / (60 * 60)) - (rDay * 24);
        var rMin = addZero(Math.floor(rTime / (60)) - (rDay * 24 * 60) - (rHour * 60));
        var rSec = addZero(Math.floor(rTime) - (rDay * 24 * 60 * 60) - (rHour * 60 * 60) - (rMin * 60));
        var rSec2 = addZero(Math.floor((rTime-Math.floor(rTime))*100));
        rDay = rDay ? rDay + '日' : '';
        rHour = rHour ? rHour + '時間' : '';
        rMin = rMin !== '00' ? rMin + '分' : '';
        timelimit.textContent = rDay + rHour + rMin + rSec + '.'+rSec2+'秒';
      }else{
        timelimit.textContent = '期限が過ぎました';
      }
    }
    gDate(rTime);
  },10);
}