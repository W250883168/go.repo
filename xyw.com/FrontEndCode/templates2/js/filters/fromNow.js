'use strict';

/* Filters */
// need load the moment.js to use this filter. 
angular.module('app')
  .filter('fromNow', function() {
    return function(date) {
      return moment(date).fromNow();
    }
  }).filter("FormatTime",function(){
  	//    秒转时间  100  >>  1分40秒    ； 0   》》  0秒
  	return function(time){
  		if(!time){time = 0}else{time = parseInt(time)};
  		if(time == 0)return 0 + "秒";
  		var h = Math.floor(time / 3600);
  		var m = (Math.floor(time / 60)) % 60;
  		var s = time % 60;
  		h　?　h = (h + "小时") : h = "";
  		m　?　m = (m + "分") : m = "";
  		s　?　s = (s + "秒") : s = "";
  		return h + m + s;
  	}
  }).filter("timeToArray",function(){
  	//   时间转数组  2016-01-01 12:12:00
  	return function(time){
  		if(time.indexOf(":")){
  			var ymd = time.substring(0,time.indexOf(" "));
  			var hms = time.substring(time.indexOf(" ")+1,time.length);
  			ymd = ymd.split("-");
  			hms = hms.split(":");
  			for(var i = 0; i < hms.length; i++){
  				ymd.push(hms[i]);
  			}
  			return ymd;
  		}
  	}
  }).filter("filterNev",function(){
  	//   导航树筛选
  	return function(items,item){
  		if(items.length > 0){
  			var tem = [];
  			for(var i = 0; i < items.length; i++){
  				for(var temname in item){
  					if(items[i][temname] == item[temname]){
  						tem.push(items[i]);
  					}
  				}
  			}
  		}else{
  			return items;
  		}
  		return tem;
  	}
  });