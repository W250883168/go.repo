'use strict';

/*    学科管理     */
app.controller("jcsjXkglContr", ['$scope', '$state','httpService', '$modal','$timeout','toaster',function($scope, $state,httpService,$modal,$timeout,toaster) {
//	//$state.go("app.qxgl",false);
  	console.log("学科管理");

	$scope.xkgl_data = {}
    $scope.xkgl_data_tree = [];

    $scope.xkglRightData = {
    	"item":{},
    	"list":[]
    };
    //   page
	$scope.page = {
		//   超始页
		"index":1,
		//   每页显示
		"oneSize":5,
		//   页码显示条数
		"pageNumber":5
	}
    //   树 已展开 保存数组
    var treeIsOpen = [];



    //    getsubjectclass 获取全部一级学科、学科查询
    var subjectclasslist = function(PageIndex,PageSize){
    	//   是展开时加入展开数组
    	treeIsOpen = [];
    	dg_tree($scope.xkgl_data_tree,function(item){
			if(item.expanded == true && item.Subjectcode != ""){
				treeIsOpen.push({"Subjectcode":item.Subjectcode,"expanded":true});
			}
		});
		Number(PageIndex) > 0 ? PageIndex = Number(PageIndex) : PageIndex = 1;
		Number(PageSize) > 0 ? PageSize = Number(PageSize) : PageSize = 15;
        var url = config.HttpUrl+"/system/us/subjectclasslist";
        var data={
        	"Usersid":config.GetUser().Usersid,
        	"Rolestype":config.GetUser().Rolestype,
        	"Token":config.GetUser().Token,
        	"Os":"WEB",
        	"PageIndex": -1
			//"PageSize":PageSize
        };
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("一级学科查询",data)
            if(data.Rcode=="1000"){
            	$scope.xkgl_data = data.Result.PageData;
            	$scope.xkgl_data_tree = outTree($scope.xkgl_data);
            	console.log($scope.xkgl_data_tree);
            	$scope.xkgl_data_tree[0].expanded = true;
				//	  分页
				//$scope.backPage = pageFn($scope.page, $scope.page.pageNumber);
				//	默认展开
				var objitem =$scope.xkglRightData.item;
				if(objitem == null || angular.equals({},objitem)){
					$scope.xkgl_data_tree[0].selected = true;
					$scope.xkglRightData.item = $scope.xkgl_data_tree[0];
					$scope.xkglRightData.list = $scope.xkgl_data_tree[0].children;
					//console.log()
				}
				//	刷新右边列表
				dg_tree($scope.xkgl_data_tree,function(item){
					if($scope.xkglRightData.item.Subjectcode == item.Subjectcode){
						$scope.xkglRightData.list = item.children;
						item.selected = true;
					}
				});
           		 //   加入展开
            	for(var a = 0; a < treeIsOpen.length; a++){
					dg_tree($scope.xkgl_data_tree,function(item){
						if(item.Subjectcode == treeIsOpen[a].Subjectcode){
							item.expanded = true;
						}
					});
				}
            }else if(data.Rcode=="1002"){
            	$scope.xkgl_data = [];
            	$scope.xkgl_data_tree = [{label:'全部学科门类',Subjectcode:"",children:[],level:0}];
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };


    //  递归树
	var dg_tree = function(tree, fn) {
		for(var i = 0; i < tree.length; i++) {
			fn(tree[i]);
//			console.log(tree[i],tree[i].children)
			if(tree[i].children.length > 0) {
				dg_tree(tree[i].children, fn);
			}
		}
	}

    /*    生成校区楼栋教室树      */
	var outTree = function(det){
		if(det.length < 1)return [{label:'全部学科门类',Subjectcode:"",children:[],level:0}];
		//
		var tree = {};
		//
		for(var a in det){
			var item = det[a];
			item.label = item.Subjectname;

			if(!tree[item.Subjectcode]) {
				tree[item.Subjectcode] = {};
			}
			//   压入当前数据
			tree[item.Subjectcode] = $.extend({},tree[item.Subjectcode],item);
			//   添加树属性
			if(!("children" in tree[item.Subjectcode])) tree[item.Subjectcode].children = [];
			//   找
			if(tree[item.Superiorsubjectcode]){
				tree[item.Superiorsubjectcode].children.push(tree[item.Subjectcode]);
			}else{
				tree[item.Superiorsubjectcode] = {
					children: [tree[item.Subjectcode]]
				};
			}
		}

		//
		var data_tree = [{label:'全部学科门类',Subjectcode:"",children:tree[""].children,level:0}];

		//   添加层级
		dg_tree(data_tree,function(item){
			for(var t_item in item.children){
				item.children[t_item].level = item.level + 1;
			}
		});

		//console.log('level',data_tree);

    	return data_tree;
	}
	//////////////////////////////////////////////////////////////////////

    //    删除学科
    var subjectclassdel = function(Subjectcode){
		if(Subjectcode=="")return false;
        var url = config.HttpUrl+"/system/us/subjectclassdel";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Subjectcode":Subjectcode
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("删除学科",data)
            if(data.Rcode=="1000"){
            	subjectclasslist();
            	toaster.pop('success', '删除成功！');
            }else{
            	toaster.pop('warning', data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };

    //添加按钮功能
    $scope.openModalAdd = function (str,active,item) {
		if(!("level" in active)){
			//alert("请选择选项！");
			$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'请选择选项！'};
                }
              }
            });
			return false;
		}

        var modalInstance = $modal.open({
            templateUrl: '../project/jcsj/html/jcsj/xkgl/modal_add.html',
            controller: 'modalxkglContr',
            windowClass: 'm-modal-xkgl',
            resolve: {
                items: function () {
                    return {'str':str,"active":active,"item":item};
                }
            }
        });

        modalInstance.result.then(function(bul) {
			console.log("刷新列表555",bul)
			if(bul){
				//   刷新列表
				//if($scope.xkglRightData.item.level!=null){
		    			subjectclasslist();
				//}

			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
    }

    //    树点击
    $scope.xkgl_tree_handler = function(branch) {
    	console.log(branch);
    	$scope.xkglRightData.item = branch;
		$scope.xkglRightData.list = branch.children;
		//	刷新右边列表
		dg_tree($scope.xkgl_data_tree,function(item){
			if(branch.Subjectcode == item.Subjectcode){
				item.selected = true;
			}else{
				item.selected = false;
			}
		});
    };


    //   删除
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'你确定要删除吗？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        subjectclassdel(item.Subjectcode);
      }
    });
  }

	/*  -------------------- 分页、页码  -----------------------  */
	$scope.backPage = {};
	/*----------------
	 //    分页对象添加页码
	 //    return  obj  分页对象
	 //    pagedata:obj  分页对象
	 //    maxpagenumber:int  显示页码数默认5个页码
	 ------------------*/
	var pageFn = function(pagedata,maxpagenumber){
		if(pagedata.length < 1)return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber){
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1 ){
				Snumber = 1
			}else if((pagedata.PageIndex + mid) > pagedata.PageCount){
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			}else{
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++){
				nub.push(Snumber + i);
			}
		}else{
			for(var i = 0; i < pagedata.PageCount; i++){
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	}

	//  翻页
	$scope.pageClick = function(pageindex){
		if(!(Number(pageindex) > 0))return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount){
			subjectclasslist(pageindex,$scope.page.oneSize);
		}
	}
	/*  -------------------- 分页、页码  -----------------------  */


    $scope.run = function(){
		subjectclasslist();
	}
	$scope.run();
}]);


app.controller("modalxkglContr",['$scope','httpService','$modal','items','$modalInstance','formValidate','toaster',function($scope,httpService,$modal,items,$modalInstance,formValidate,toaster){
	console.log("基础数据-学科管理-弹窗");
	$scope.items=items;
	console.log("弹窗",items);

	$scope.form = {
		//一级学科
		"subject":{}
	}
	$scope.form.subject = {
		"Subjectcode":'',
		"Subjectname":"",
		"Superiorsubjectcode":""
	}


	//    验证学科代码
    $scope.changeSubjectcode = function(){
    	//   本级学科增加长度
    	var xk_lang = 4;
    	//   学科代码
    	var code = $scope.form.subject.Subjectcode;
    	//   上级学科代码
    	var pcode = $scope.items.active.Subjectcode;
    	//    验证输入必须为数字
    	if(!(/^[0-9]*$/.test(code))){

    		$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'只能输入上级代码加4位的数字！'};
                }
              }
            });
    		//alert('只能输入上级代码加4位的数字！');
    		code = pcode;
    	}

    	//
    	if(code.length < pcode.length){
			//   里面没有上级学科代码  则 放入上级学科代码
			code = pcode;
    	}else{
    		//   有  验证前面的数字是不是上级学科代码
    		if(code.substr(0,pcode.length) == pcode){
    			//   是上级学科代码 。限制输入长度
    			if(code.length > pcode.length + xk_lang){
    				code = code.substr(0,code.length - 1);
    			}else{
    				console.log(code);
    			}
    		}else{
    			//   不是上级学科代码 。 放入上级学科代码
    			code = pcode;
    		}
    	}
    	//
    	$scope.form.subject.Subjectcode = code;
    }


	//   一级学科添加
    var subjectclassadd = function(){
    	//   验证CODE
    	if(!(formValidate($scope.form.subject.Subjectcode).isNumber().outMsg(2800).isOk))return false;
    	if($scope.items.active.Subjectcode == $scope.form.subject.Subjectcode){
    		$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'学科代码不能与上级学科代码相同！'};
                }
              }
            });
            return false;
    	}


		$scope.changeSubjectcode();

		var url = config.HttpUrl+"/system/us/subjectclassadd";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Subjectcode":$scope.form.subject.Subjectcode,
            "Subjectname":$scope.form.subject.Subjectname,
            "Superiorsubjectcode":$scope.form.subject.Superiorsubjectcode
		};

        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("一级学科添加",data)
            if(data.Rcode=="1000"){
            	toaster.pop('success', '添加成功！');
            	$modalInstance.close(true);
            }else{
            	toaster.pop('warning', data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };

    //   一级学科修改
    var subjectclasschange = function(){
		//   验证CODE
    	if(!(formValidate($scope.form.subject.Subjectcode).isNumber().outMsg(2800).isOk))return false;
    	if($scope.items.active.Subjectcode == $scope.form.subject.Subjectcode){
    		$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'学科代码不能与上级学科代码相同！'};
                }
              }
            });
            return false;
    	}

		$scope.changeSubjectcode();

        var url = config.HttpUrl+"/system/us/subjectclasschange";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Subjectcode":$scope.form.subject.Subjectcode,
            "Subjectname":$scope.form.subject.Subjectname,
            "Superiorsubjectcode":$scope.form.subject.Superiorsubjectcode,
            "Id":$scope.form.subject.Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("学科修改",data)
            if(data.Rcode=="1000"){
            	toaster.pop('success', '修改成功！');
            	$modalInstance.close(true);
            }else{
            	toaster.pop('warning', data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };

    $scope.ok = function(){
		if($scope.items.str == "add"){
			if($scope.items.active.level!=null){
					subjectclassadd();
			}
		}
		if($scope.items.str == "edit"){
			if($scope.items.active.level!=null){
					subjectclasschange();
			}
		}
	}

    $scope.run = function(){
		if($scope.items.str == "add"){
			switch($scope.items.active.level){
				case 0:
					$scope.form.subject.Superiorsubjectcode = "";
					$scope.form.subject.Subjectcode = $scope.items.active.Subjectcode;
				break;
				case 1:
		            $scope.form.subject.Superiorsubjectcode = $scope.items.active.Subjectcode;
		            $scope.form.subject.Subjectcode = $scope.items.active.Subjectcode;
				break;
				case 2:
					$scope.form.subject.Superiorsubjectcode = $scope.items.active.Subjectcode;
					$scope.form.subject.Subjectcode = $scope.items.active.Subjectcode;
				break;
			}
		}
		if($scope.items.str == "edit"){
			switch($scope.items.active.level){
				case 0:
					$scope.form.subject.Subjectcode = $scope.items.item.Subjectcode;
					$scope.form.subject.Subjectname = $scope.items.item.Subjectname;
					$scope.form.subject.Superiorsubjectcode = "";
					//$scope.form.subject = $.extend({}, $scope.form.subject, $scope.items.item);
					$scope.form.subject.Id = $scope.items.item.Id;
				break;
				case 1:
		            $scope.form.subject = $.extend({}, $scope.form.subject, $scope.items.item);
		            $scope.form.subject.Id = $scope.items.item.Id;
				break;
				case 2:
					$scope.form.subject = $.extend({}, $scope.form.subject, $scope.items.item);
		            $scope.form.subject.Id = $scope.items.item.Id;
				break;
			}
		}
	}
	$scope.run();
	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}
}]);
