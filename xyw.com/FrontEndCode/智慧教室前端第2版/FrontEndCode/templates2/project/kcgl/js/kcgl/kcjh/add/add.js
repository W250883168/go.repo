'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理-课程计划-添加
 */

app.filter('propsFilter', function() {
	return function(items, props) {
		var out = [];
		if(angular.isArray(items)) {
			items.forEach(function(item) {
				var itemMatches = false;
				var keys = Object.keys(props);
				for(var i = 0; i < keys.length; i++) {
					var prop = keys[i];
					try {
						if(typeof(props[prop]) == "number") {
							var text = props[prop];
							if(item[prop].toString() == text.toString()) {
								itemMatches = true;
								break;
							}
						} else {
							var text = props[prop].toString().toLowerCase();
							if(item[prop].toString().toLowerCase().indexOf(text) !== -1) {
								itemMatches = true;
								break;
							}
						}
					} catch(e) {}
				}
				if(itemMatches) {
					out.push(item);
				}
			});
		} else {
			out = items;
		}
		return out;
	};
});

/*    课程管理-课程计划-添加       */
app.controller("kcglKcjhAddContr", ['$scope', 'httpService','$modal', function($scope,httpService,$modal) {
	console.log("课程管理-课程计划-添加")
	
	
	//   课程信息
	$scope.form = {
		//   课程信息
		"lesson":{},
		//   章节信息
		"section":[],
		//  向下级页面传递的操作状态 
		"operation":"add"
	}
	
	//   课程信息
	$scope.form.lesson = {
		//     所属学科 ID
		"scienceCode":"",
		//     所属学科 名
		"scienceName":"",
		//    学科目录
		"scienceNameTree":"",
		//   课程ID
		"courseId": "",
		//   课程名称
		"courseName": "",
		//   课程图标
		"pic":null,
		//    课程类型 1普通课2公开课
		"courseType":"1",
		//   上课老师
		"teacher":[],
		//   上课老师显示文本
		"teacherText":"",
		//    班级
		"class":[],
		//    班级显示文本
		"classText":"",
		//    课程详情
		"courseDetails":"",
		//    总章节数
		"allChapter":"0"
	}
	
	//    打开弹窗  选择学科
	$scope.modalOpenScience = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_science.html',
			controller: 'modalGetScienceCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectItem) {
			console.log(selectItem)
			if(!selectItem){
				$scope.form.lesson.scienceCode = "";
				$scope.form.lesson.scienceName = "";
			}else{
				$scope.form.lesson.scienceCode = selectItem.Subjectcode;
				$scope.form.lesson.scienceName = selectItem.Subjectname;
				$scope.form.lesson.scienceNameTree = selectItem.SubjectnameTree;
			}
			console.log($scope.form)
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	
	//  上传图片
	var handleFileSelect = function(evt) {
		var file = evt.currentTarget.files[0];
		if(!file && !/image\/\w+/.test(file.type))return false;
		var reader = new FileReader();
		reader.onload = function(evt) {
			$scope.$apply(function($scope) {
				$scope.form.lesson.pic = evt.target.result;
			});
		};
		reader.readAsDataURL(file);
	};
	$(document).on('change', '#fileInput_add', handleFileSelect);
	
	//    清除图片
	$scope.closePic = function(){
		$scope.form.lesson.pic = "";
	}
	
	//   去重复
	//   oldsArr  数组
	//   oldsId   属性名
	//   
	//   返回         数组
	var removeRepeat = function(oldsArr,oldsId,newsObj,newsId){
		var temp = false;
		if(oldsArr.length == 0 && oldsArr[oldsId] != ""){
			oldsArr.push(newsObj);
			return oldsArr;
		}
		for(var i = 0; i < oldsArr.length; i++){
			if(newsObj[newsId] == oldsArr[i][oldsId]){
				temp = false;
				break;
			}else{
				temp = true;
			}
		}
		if(temp){
			oldsArr.push(newsObj);
			return oldsArr;
		}else{
			return oldsArr;
		}
	}
	
	
	//   上课老师
	var getTeacherListText = function(){
		var teacherList = [];
		var teacherText = "";
		for(var i = 0; i < $scope.form.section.length; i++){
			for(var b = 0; b < $scope.form.section[i].classRoom.length; b++){
				if($scope.form.section[i].classRoom[b].teacherId != ""){
					teacherList = removeRepeat(teacherList,'teacherId',$scope.form.section[i].classRoom[b],'teacherId');
				}
			}
		}
		for(var i = 0; i < teacherList.length; i++){
			teacherText += teacherList[i].teacherName + ",";
		}
		teacherText = teacherText.substring(0,teacherText.length - 1);
		//
		$scope.form.lesson.teacher = teacherList;
		return teacherText;
	}
	
	//   上课班级
	var getClassListText = function(){
		var classList = [];
		var classText = "";
		for(var i = 0; i < $scope.form.section.length; i++){
			for(var b = 0; b < $scope.form.section[i].classRoom.length; b++){
				if($scope.form.section[i].classRoom[b].classId != ""){
					classList = removeRepeat(classList,'classId',$scope.form.section[i].classRoom[b],'classId');
				}
			}
		}
		for(var i = 0; i < classList.length; i++){
			classText += classList[i].className + ",";
		}
		classText = classText.substring(0,classText.length - 1);
		//
		$scope.form.lesson.class = classList;
		return classText;
	}
	
	//   总章节数
	var getSectionListText = function(){
		$scope.form.lesson.allChapter = $scope.form.section.length;
		console.log($scope.form.lesson.allChapter)
	} 
	
	//    打开弹窗  添加章节
	var modalOpenAdd = function(item) {
		var modalInstance = $modal.open({
			templateUrl: '../project/kcgl/html/kcgl/modal/modal_add_add.html',
			controller: 'kcglKcjhAddAddContr',
			windowClass: 'm-modal-kcgl-add',
			resolve: {
				items: function() {
					return item;
				}
			}
		});

		modalInstance.result.then(function(form) {
			console.log(form)
			if(!form){
				//$scope.deviceText = "";
			}else{
				$scope.form.section.push(form);
				//  回写老师
				getTeacherListText();
				//  班级
				getClassListText();
				//   章节数
				getSectionListText();
				//   上课老师
				$scope.form.lesson.teacherText = getTeacherListText();
				//   上课班级
				$scope.form.lesson.classText = getClassListText();
			}
			console.log($scope.form)
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	//   操作
	$scope.operationClick = function(str,item){
		switch(str){
			case "add":
				$scope.form.operation = "add";
				modalOpenAdd(item);
			break;
			case "details":
				$scope.form.operation = "details";
				modalOpenAdd(item);
			break;
			case "edit":
				$scope.form.operation = "edit";
				modalOpenAdd(item);
			break;
			case "remove":
			
			break;
		}
	}
	
	
	$scope.ok = function() {
		//console.log(item)
		//   还差验证
		//$modalInstance.close();
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
	
	//  run
	var run = function(){
		
	}
	run();
	
}]);