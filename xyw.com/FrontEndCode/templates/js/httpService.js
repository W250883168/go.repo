/**
 * Created by Administrator on 2016/7/2.
 */
app.factory('httpService', function ($http, $q) {
    return {

        //========================================================================
        // 下面是通过$http访问后台进行记录的增、删、改、查
        //========================================================================

        //查询
        ajaxGet: function (url, params) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "GET",
                params: params
            })
                .success(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("success");
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("Service error");
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        }
        ,

        //删除
        ajaxDelete: function (url) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "DELETE",
                params: {}
            })
                .success(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("success");
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("Service error");
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        }
        ,

        //增加(POST)
        ajaxPost: function (url, data) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "POST",
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                data: data,
                params: {}

            })
                .success(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("success");
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("Service error");
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        },

        //修改
        ajaxPut: function (url, row) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "PUT",
                data: row,
                params: {}
            })
                .success(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("success");
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    Metronic.unblockUI();
                    //alert("Service error");
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        }


    };
});

