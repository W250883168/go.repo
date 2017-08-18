'use strict';
/**
 * Created by Administrator on 2016/7/21.
 */
angular.module('app').factory('httpService', function ($http, $q) {
    return {

        //========================================================================
        // 下面是通过$http访问后台进行记录的增、删、改、查
        //========================================================================

        //GET请求
        ajaxGet: function (url, params) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "GET",
                params: params
            })
                .success(function (data, status, headers, config) {
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        },
        //DELETE请求
        ajaxDelete: function (url) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "DELETE",
                params: {}
            })
                .success(function (data, status, headers, config) {
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        },
        //POST请求
        ajaxPost: function (url, data, out) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "POST",
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                data: data,
                timeout:out,
                params: {}
            })
                .success(function (data, status, headers, config) {
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        },
        //PUT请求
        ajaxPut: function (url, row) {
            var deferred = $q.defer();
            $http({
                url: url,
                method: "PUT",
                data: row,
                params: {}
            })
                .success(function (data, status, headers, config) {
                    deferred.resolve(data);
                })
                .error(function (data, status, headers, config) {
                    deferred.reject('Service error ');
                })
            return deferred.promise;
        }
    };
});

