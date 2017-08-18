'use strict';

/**
 * @function 表单验证服务
 * @version 1.0
 * @author wu
 * @return {}
 */
angular.module('app').factory('formValidate', ['MsgText','$modal',function (MsgText,$modal) {

	//console.log($translate);

	function formThis(data){

		var _this = formThis;

		//    正则字符串
		var RegExpStr = {
			//正则
		   reg_email: /^\w+\@[a-zA-Z0-9]+\.[a-zA-Z]{2,4}$/i, //验证邮箱
		   reg_num: /^\d+$/,         //验证数字
		   reg_chinese: /^[\u4E00-\u9FA5]+$/,     //验证中文
		   reg_mobile: /^1[34578]{1}[0-9]{9}$/,    //验证手机
		   reg_node: /^\w{16}$/i, //验证节点编号
		   reg_idcard: /^\d{14}\d{3}?\w$/      //验证身份证
		}

		//   msg 提示文本。
		var  msg = "";
		//   提示类型 提示
		var vType = 'alert';
		//   返回data
		var Rdata = data;
		//   验证是否通过      true:验证通过，false:验证不通过
		var isOk = false;

		//   输出提示   str:2400   or   str:"str"
		_this.outMsg = function(str){
			console.log(angular.isNumber(str));
			if(str && angular.isNumber(str)){
				str = MsgText.alert.verify[str];
			}else if(str == undefined){
				str = _this.msg;
			}else{
				str = str;
			}
			//console.log(str);
			if(_this.isOk == false){
				switch(vType){
					case "alert":
						//
						//alert(str);
            var modalInstance = $modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'info',"msg":str};
                }
              }
            });
						return _this;
					break;
					case "error":
						//
						// alert(str);
            var modalInstance = $modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":str};
                }
              }
            });
						return _this;
					break;
					case "ask":
						//
						//alert(str);
            var modalInstance = $modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'warning',"msg":str};
                }
              }
            });
						return _this;
					break;
					default:
						//
						//alert(str);
            var modalInstance = $modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'info',"msg":str};
                }
              }
            });
						return _this;
				}
			}else{
				return _this;
			}

		}

		//   最大长度 i:int
		_this.maxLength = function(i){
			_this.vType = 'alert';
			if(data && data.length < i){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2300].replace(/{i}/, i);
				return _this;
			}
		}
		//   最小长度  i:int
		_this.minLength = function(i){
			_this.vType = 'alert';
			if(data && data.length > i){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2301].replace(/{i}/, i);
				return _this;
			}
		}
		//   验证邮箱
		_this.isEmail = function(){
			_this.vType = 'alert';
			if(data && RegExpStr.reg_email.test(data)){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2201];
				return _this;
			}
		}
		
		//   验证手机
		_this.isMobile = function(){
			_this.vType = 'alert';
			if(data && RegExpStr.reg_mobile.test(data)){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2205];
				return _this;
			}
		}

		//   验证数字
		_this.isNumber = function(){
			_this.vType = 'alert';
			if(data && RegExpStr.reg_num.test(data)){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2203];
				return _this;
			}
		}

		//   验证节点编号
		_this.isNode = function(){
			_this.vType = 'alert';
			if(data && RegExpStr.reg_node.test(data)){
				_this.Rdata = data;
				_this.msg = "";
				_this.isOk = true;
				return _this;
			}else{
				_this.isOk = false;
				_this.msg = MsgText.alert.verify[2204];
				return _this;
			}
		}


		return _this;

	}
	return formThis;
}]);


//     msg 服务
app.factory('MsgText', function() {
	return {
		//   错误
		"error":{
			"1000":"错误！"
		},
		//   提示
		"alert":{
			"2000":"输入错误！",
			//   验证信息
			"verify":{
				//
				"2100":"输入不能为空！",
				//
				"2101":"输入物业编码不能为空！",
				"2102":"输入物业名称不能为空！",
				"2103":"输入所在楼层不能为空！",
				"2104":"输入单元不能为空！",
				"2105":"输入房号不能为空！",
				"2106":"输入户型不能为空！",
				"2107":"输入建筑面积不能为空！",
				"2108":"输入建筑面积只能为数字！",
				"2109":"输入租凭面积不能为空！",
				"2110":"输入租凭面积只能为数字！",
				"2111":"输入产权证号不能为空！",
				"2112":"输入购置日期不能为空！",
				"2113":"输入使用日期不能为空！",
				"2114":"输入详细位置不能为空！",
				"2115":"输入备注不能为空！",
				"2116":"选择是否封存不能为空！",
				"2117":"选择是否删除不能为空！",
				//
				"2200":"邮箱不能为空！",
				"2201":"邮箱填写错误！",
				"2202":"数字不能为空！",
				"2203":"数字格式错误！",
				"2204":"节点编号错误必须16位字母与数字组合！",
				"2205":"手机号码格式错误！",
				//
				"2300":"输入最大文本长度为{i}！",
				"2301":"输入最小文本长度为{i}！",

				//  24节点配置验证提示文本
				"2400":"节点型号名称不能为空！",
				//
				"2401":"命令代码不能为空!",
				"2402":"节点型号Id不能为空!",
				"2403":"命令名称不能为空!",
				"2404":"请求类型不能为空，且只能为post/get/put/delete四种类型!",
				"2405":"请求地址不能为空!",
				"2406":"请求地址参数不能为空!",
				//
				"2407":"节点编号不能为空!",
				"2408":"节点型号不能为空!!",
				"2409":"教室位置不能为空",
				"2410":"来源路由Ip地址不能为空!",
				"2411":"最新上报数据时间不能为空!",
				//	25设备配置验证提示文本
				"2500":"设备型号Id不能为空!",
				"2501":"设备型号名称不能为空!",
				"2502":"设备型号类型不能为空!",
				//
				"2503":"设备型号控制命令代码不能为空!",
				"2504":"设备型号控制命令名称不能为空!",
				"2505":"设备型号请求地址不能为空!",
				"2506":"设备型号URI查询参数不能为空!",
				"2507":"设备型号负载不能为空!",
				"2508":"设备型号请求类型不能为空!",
				//
				"2509":"设备型号状态编码不能为空!",
				"2510":"设备型号状态名称不能为空!",
				"2511":"设备型号状态序号不能为空!",
				"2512":"设备型号命令负载不能为空!",
				"2513":"设备型号状态值匹配串不能为空!",
				"2514":"设备型号状态命令编码值不能为空!",
				"2515":"设备型号状态命令名称值不能为空!",
				//
				"2516":"设备型号设备名称不能为空!",
				"2517":"设备型号资产编号不能为空!",
				"2518":"设备型号设备品牌不能为空!",
				"2519":"设备型号所在教室不能为空!",
				"2523":"节点编号或接入方式节点编号必填一项!",
				"2524":"节点编号与接入方式节点编号必须相同!",
				//
				"2521":"故障现象常用词条名称不能为空!",
				//
				"2522":"设备型号故障分类名称不能为空!",
				//"2523":"节点编号或接入方式节点编号必填一项!",
				//"2524":"节点编号与接入方式节点编号必须相同!",

				//	26 设备故障验证提示文本
				//"2600":"故障ID不能为空!",
				"2601":"故障设备不能为空，请选择设备!",
				"2602":"故障现象不能为空，请输入故障现象!",
				"2603":"故障发生时间不能为空，请选择故障发生时间!",
				"2604":"请选择故障设备是否可用!",
				"2605":"维修人不能为空!",
				"2606":"维修人电话不能为空!",
				"2607":"维修完成时间不能为空!",
				"2608":"维修设备是否可用不能为空!",
				"2609":"维修结果不能为空!",
				//  27 联动场景验证提示文本
				"2700":"联动场景名称不能为空!",
				"2701":"联动场景触发条件不能为空!",
				"2702":"联动场景响应事件不能为空!",
				"2703":"联动场景位置不能为空!",
				// 28 基础数据验证提示文本
				"2800":"一级学科代码必须是4位以内整数!",
				"2801":"二级学科代码必须是上一级的2位整数再加上2位整数!",
				"2802":"三级学科代码必须是上一级的4位整数再加上2位整数!",
				//  校区管理
				"2803":"校区代码不能为空!",
				"2804":"校区名称不能为空!",
				"2805":"楼栋代码不能为空!",
				"2806":"楼栋名称不能为空!",
				"2807":"楼层代码不能为空!",
				"2808":"楼层名称不能为空!",
				"2809":"教室代码不能为空!",
				"2810":"教室名称不能为空!",
				//  学院管理
				"2811":"学院代码不能为空!",
				"2812":"学院名称不能为空!",
				"2813":"系列代码不能为空!",
				"2814":"系列名称不能为空!",
				"2815":"班级代码不能为空!",
				"2816":"班级名称不能为空!",
				"2817":"入学年份不能为空!",
				// 29 课程管理验证提示文本         基础课程
				"2900":"章节名称不能为空!",
				//   课程计划
				"2901":"课程名称不能为空!",
				"2902":"班级不能为空!",
				"2903":"上课老师不能为空!",
				"2904":"章节配置章节名称不能为空!",
				"2905":"老师不能为空!",
				//	基础数据
				"3101":"用户账号不能为空!",
				"3102":"账号密码不能为空!",
				"3103":"用户角色不能为空!",
				"3104":"真实姓名不能为空!",
				"3105":"属所学院不能为空!",
				"3106":"昵称不能为空!",
				//	基础数据-学生管理
				"3201":"入学年份不能为空!",
				"3202":"家庭地址不能为空!",
				"3203":"所在班级不能为空!"
			}
		},
		//   问
		"ask":{
			"3000":"你确定要删除吗？"
		}
	};
});

/**
 * 筛选出栏目/模块下的所有子栏目
 * ndata 栏目数组;navid要筛选的栏目ID
 * @return [];栏目/模块数组
 */
//app.factory('getNavSon', function() {
//		return function(){
//
//		}
//}
