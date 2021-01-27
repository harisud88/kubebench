/*
 * htmlson.js v.1 (Adalen VLADI, Redjan Ymeraj) | MIT
 * Github: https://github.com/adalenv/htmlson.js
 */
/***** Helper methods *****/
function isObject(value) {
    return value && typeof value === 'object' && value.constructor === Object;
}

function getDepth(obj) {
    var depth = 0;
    if (obj.children) {
        obj.children.forEach(function (d) {
            var tmpDepth = getDepth(d);

            if (tmpDepth > depth) {
                depth = tmpDepth
            }
        })
    }
    return 1 + depth
}
/***** Helper methods *****/
(function ($) {
    $.fn.htmlson = function (configs) {
        var scope = this;
        var autoHeaderKeys = [];
        var thead = '';
        var tbody = '';

        /***** Start parse configurations *****/
        if (typeof configs.data !== "object") {
            console.error("htmlson.js Error: No data passed!");
            return;
        }

        if (typeof configs.headers !== "object") {
            configs.headers = {};
        }

        if (typeof configs.debug !== "boolean") {
            configs.debug = false;
        }
        /***** End parse configurations *****/

        function initialize () {
            // validate the
            if($.isEmptyObject(configs.data)){
                // set empty table
                scope.html('');
                return;
            }

            /***** Start set headers *****/
            thead = '<thead>';

            autoHeaderKeys = Object.keys(configs.data[0]);

            for (var i = 0; i < autoHeaderKeys.length; i++) {
                if (configs.headers[i] === undefined) {
                    thead += '<th style="background-color: #007bff;color: #fff;">'+autoHeaderKeys[i]+'</th>';//if auto header
                } else {
                    thead += '<th style="background-color: #007bff;color: #fff;">'+configs.headers[i]+'</th>';//if user defined header
                }
            }

            thead += '</thead>';
            /***** End set headers *****/

            /***** Start set body *****/
            tbody = '<tbody>';

            for (var i in configs.data) {
                tbody += '<tr class="table-active">';

                var array = $.map(configs.data[i], function (value, index) {
                    return value;
                });

                for (var j in array) {
                    if (!isObject(array[j])) {                      //if not object
                        tbody += '<td><a id="download" href="#" onclick="myfunc(event)">'+array[j]+'<a></td>'
                    } else {                                        //if object convert to ul
                        tbody += '<td><ul>';
                        var ob = $.map(array[j], function (value, index) {
                            return value;
                        });
                        for (var h in ob) {
                            tbody += '<li>'+ob[h]+'</li>';
                        }
                        tbody += '</ul></td>';
                    }
                }

                tbody += '</tr>';
            }

            tbody += '</tbody>';

            /***** End set body *****/

            /***** Start generate output *****/
            scope.html(thead + tbody);
            /***** End generate output *****/
        }

        initialize();

        /***** Start debug *****/
        if (configs.debug) {
            var log = function (l) {
                console.log(l);
            };
            log('Debug: true');
            log('Object: ' + JSON.stringify(configs.data));
            log('Object depth: ' + getDepth(configs.data));
            log('Auto headers: ' + JSON.stringify(autoHeaderKeys));
            log('Headers set: ' + JSON.stringify(configs.headers));
            log('Table head: ' + thead);
            log('Table body: ' + tbody);
        }
        /***** End debug *****/


        return scope;
    };
}(jQuery));

(function ($) {
    $.fn.htmlcustomson = function (configs) {
        var scope = this;
        var autoHeaderKeys = [];
        var thead = '';
        var tbody = '';

        /***** Start parse configurations *****/
        if (typeof configs.data !== "object") {
            console.error("htmlson.js Error: No data passed!");
            return;
        }

        if (typeof configs.headers !== "object") {
            configs.headers = {};
        }

        if (typeof configs.debug !== "boolean") {
            configs.debug = false;
        }
        /***** End parse configurations *****/

        function initialize () {
            // validate the
            if($.isEmptyObject(configs.data)){
                // set empty table
                scope.html('');
                return;
            }

            /***** Start set headers *****/
            thead = '<thead>';

            autoHeaderKeys = Object.keys(configs.data[0]);

            for (var i = 0; i < autoHeaderKeys.length; i++) {
                if (configs.headers[i] === undefined) {
                    thead += '<th style="background-color: #007bff;color: #fff;">'+autoHeaderKeys[i]+'</th>';//if auto header
                } else {
                    thead += '<th style="background-color: #007bff;color: #fff;">'+configs.headers[i]+'</th>';//if user defined header
                }
            }

            thead += '</thead>';
            /***** End set headers *****/

            /***** Start set body *****/
            tbody = '<tbody>';

            for (var i in configs.data) {

                var array = $.map(configs.data[i], function (value, index) {
                    return value;
                });
                if (array[0] == "PASS")
                {
                    tbody += '<tr class="table-success">';
                }
                else{
                    tbody += '<tr class="table-danger">';
                }
                for (var j in array) {
                    if (!isObject(array[j])) { 
                        if (array[j] == "PASS" || array[j] == "FAIL")
                        {                     //if not object
                        tbody += '<th scope="row">'+array[j]+'</td>'
                        }
                        else{
                            tbody += '<td>'+array[j]+'</td>'
                        }
                    } else {                                        //if object convert to ul
                        tbody += '<td><ul>';
                        var ob = $.map(array[j], function (value, index) {
                            return value;
                        });
                        for (var h in ob) {
                            tbody += '<li>'+ob[h]+'</li>';
                        }
                        tbody += '</ul></td>';
                    }
                }

                tbody += '</tr>';
            }

            tbody += '</tbody>';

            /***** End set body *****/

            /***** Start generate output *****/
            scope.html(thead + tbody);
            /***** End generate output *****/
        }

        initialize();

        /***** Start debug *****/
        if (configs.debug) {
            var log = function (l) {
                console.log(l);
            };
            log('Debug: true');
            log('Object: ' + JSON.stringify(configs.data));
            log('Object depth: ' + getDepth(configs.data));
            log('Auto headers: ' + JSON.stringify(autoHeaderKeys));
            log('Headers set: ' + JSON.stringify(configs.headers));
            log('Table head: ' + thead);
            log('Table body: ' + tbody);
        }
        /***** End debug *****/


        return scope;
    };
}(jQuery));