/**
 * Created by Administrator on 2016/7/2.
 */
document.cookie='LoginUser={"Loginuser":"admin","Loginpwd":"","Usersid":1,"Rolestype":1,"Truename":"刘某某","Nickname":"刘管理","Userheadimg":"/web/upfile/20170407/1491536854441839900.jpg","Userphone":"13788889977","Userstate":0,"Usersex":1,"Usermac":"","Birthday":"1988-08-10","Token":"admin|1|1","ThirdPartyId":"","Os":""}';

//alert(document.cookie);
function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg)) {
        return unescape(arr[2]);
    } else {
        return null;
    }
}

//alert(123)
//alert(localStorage.getItem("LoginUser"));
//alert(getCookie('LoginUser'));

if (getCookie('LoginUser') == null) {
    document.cookie = 'LoginUser={"Loginuser":"","Usersid":null,"Nickname":"","Rolestype":null,"Token":"","Modelname":"设备管理"}';
}

/**
*IF IOS
*/
if (true) {
    if (true) {
        localStorage.setItem("LoginUser", getCookie('LoginUser'));
    } else {
        localStorage.setItem("LoginUser", LocalStorageMy.getItem("LoginUser"));
    }
} else {
    localStorage.setItem("LoginUser", LocalStorageMy.getItem("LoginUser"));
}

var config = {
    HttpUrl: "",
    //HttpUrl: "http://localhost:8080",
    loginuser: null,
    GetUser: function () {
        if (this.loginuser == null) {
            if (localStorage.getItem("LoginUser") != null && localStorage.getItem("LoginUser") != "undefined" && localStorage.getItem("LoginUser") != "") {
                this.loginuser = jQuery.parseJSON(localStorage.getItem("LoginUser"));
            } else {
                this.loginuser = null;
                //window.location.href = "/web/html/login.html";
            }
        }
        return this.loginuser;
    }
};


/**
 * app
 */
var app = angular.module("app", ['ui.bootstrap', 'ngTouch']);

//
//app.config(['$httpProvider', function ($httpProvider) {
//    $httpProvider.defaults.withCredentials = true;
//
//}]);
app.config(['$locationProvider', function ($locationProvider) {
    $locationProvider.html5Mode(true);
}]);



/* Setup global settings */
app.factory('settings', ['$rootScope', '$window', function ($rootScope, $window) {
    var settings = {
        coapServer: "http://192.168.0.209:8090",
        deviceImg: "/web/upfile/device/",
        devicePage: "/web/html/",
        getClassroomInfoUrl: "/basicset/getclassroominfo",
        immediatelyRefreshTime: 50, //立即刷新时间
        fixRefreshTime: 3000 //固定刷新时间
    };
    return settings;
}]);







/*    过滤器     开始       */
angular.module("app").filter("FormatTime",function(){
  	//    秒转时间  100  >>  1分40秒    ； 0   》》  0秒
  	return function(time){
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
  
/*    过滤器     结束       */


/*    指令    开始       */

//    unSlider 轮播图  指令
angular.module("app").directive('unSlider', [function () {
    function link($scope, element, attrs) {
        //监听options变化
        if (attrs.uiOptions) {
            attrs.$observe('uiOptions', function () {
                var options = $scope.$eval(attrs.uiOptions);
                if (angular.isObject(options)) {
                    element.unslider(options);
                }
            }, true);
        }
    }
    return {
        restrict: 'A',
        link: link
    };
}]);



/*    监听滑动      */
angular.module("app").directive('imageController',['$swipe',function($swipe){
    return {
        restrict:'EA',
        link:function(scope,ele,attrs,ctrl){
        	var dom_body = $("body");
        	var dom_unslider = null;
        	var start_left = 0;
        	var dom_unslider_width = 0;
            var startX,startY,locked=false;
            $swipe.bind(ele,{
                'start':function(coords){
                	dom_unslider = ele.find(".unslider-wrap");
                	dom_unslider_div = ele.find(".unslider-horizontal");
                	start_left = dom_unslider.position().left;
                	dom_unslider_width = dom_unslider.width();
                	dom_unslider_li_width = dom_unslider.find("> li").eq(0).width();
                    startX = coords.x;
                    startY = coords.y;
                },
                'move':function(coords){
                	dom_unslider.css("left",start_left + coords.x - startX);
                },
                'end':function(coords){
                	var index = 0;
					ele.find(".unslider-nav ol li").each(function(a,b){if($(b).hasClass('unslider-active')){index = $(b).attr('data-slide');}});
                	//   <--
                	if((coords.x - startX) < -80 && (coords.x - startX) < 0){
                		if(start_left + coords.x - startX < -(dom_unslider_width - dom_unslider_li_width)){
                			dom_unslider_div.unslider('animate:last');
                		}else{
                			dom_unslider_div.unslider('next');
                		}
                	//    -->
					}else if((coords.x - startX) > 80 && (coords.x - startX) > 0){
						if(start_left + coords.x - startX > 0){
							dom_unslider_div.unslider('animate:first');
						}else{
							dom_unslider_div.unslider('prev');
						}
					//    ===
					}else if(coords.x - startX != 0){
						dom_unslider_div.unslider('animate:' + index);
					}
                },
                'cancel':function(coords){
                }
            });
        }
    }
}]);







/*    指令    结束       */



/*    弹窗指令     */
angular.module("app").directive('alertBar', [function() {
	return {
		restrict: 'EA',
		templateUrl: '../html/modal/modal_alert.html',
		scope: {
			message: "=",
			type: "="
		},
		link: function(scope, element, attrs) {

			scope.hideAlert = function() {
				scope.message = null;
				scope.type = null;
			};

		}
	};
}]);
//     弹窗服务
angular.module("app").factory('TipService', ['$timeout', function($timeout) {
	return {
		message: null,
		type: null,
		setMessage: function(msg, type) {

			this.message = msg;
			this.type = type;

			//提示框显示最多3秒消失
			var _self = this;
			$timeout(function() {
				_self.clear();
			}, 3000);
		},
		clear: function() {
			this.message = null;
			this.type = null;
		}
	};
}]);
       