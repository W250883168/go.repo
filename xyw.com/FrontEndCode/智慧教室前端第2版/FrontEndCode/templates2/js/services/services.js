/**
 * ---------------------------------------------
 * 
 * (c) Copyright 2017 wu. All Rights Reserved.  
 * 
 * 后台管理系统服务
 * 
 * ---------------------------------------------
 */



//     中控制面板用弹窗服务
app.factory('TipService', ['$timeout', function($timeout) {
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


/**
 * 后台提示弹窗服务2
 * @return {} obj 
 * @return {}.msg  提示文本 
 * @return {}.type 类型 type:'alert , error , ask ' /  提示 ， 错误 ， 问
 * @return {}.title  标题
 * @return {}.isOk  返回 点的 取消：false  ， 还是 确定 ：true
 */
app.factory('alertService', ['$timeout','$modal', function($timeout,$modal) {
	return {
		msg: null,
		type: null,
		title:null,
		isOk:false,
		clear: function() {
			this.message = null;
			this.type = null;
			this.title = null;
		},
		//   type:'alert , error , ask ' /  提示 ， 错误 ， 问
		openMsg:function(msg,type,title){
			var modalInstance = $modal.open({
	            templateUrl: '../html/modal/modal_alert.html',
	            controller: 'modalAlertConter',
	            windowClass: 'm-modal-alert',
	            backdrop:"static",
	           	resolve: {
	                items: function () {
	                    return {'msg':msg,"title":title,"type":type};
	                }
	            }
	        });
	        
	        modalInstance.result.then(function(bul) {
				this.isOk = bul;
			});
		}
	};
}]);	
	
	