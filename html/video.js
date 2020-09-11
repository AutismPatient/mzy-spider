/**
 * @license
 * Video.js 7.7.5 <http://videojs.com/>
 * Copyright Brightcove, Inc. <https://www.brightcove.com/>
 * Available under Apache License Version 2.0
 * <https://github.com/videojs/video.js/blob/master/LICENSE>
 *
 * Includes vtt.js <https://github.com/mozilla/vtt.js>
 * Available under Apache License Version 2.0
 * <https://github.com/mozilla/vtt.js/blob/master/LICENSE>
 */

(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? module.exports = factory(require('global/window'), require('global/document')) :
        typeof define === 'function' && define.amd ? define(['global/window', 'global/document'], factory) :
            (global = global || self, global.videojs = factory(global.window, global.document));
}(this, function (window$1, document) { 'use strict';

    window$1 = window$1 && window$1.hasOwnProperty('default') ? window$1['default'] : window$1;
    document = document && document.hasOwnProperty('default') ? document['default'] : document;

    var version = "7.7.5";

    /**
     * @file create-logger.js
     * @module create-logger
     */

    var history = [];
    /**
     * Log messages to the console and history based on the type of message
     *
     * @private
     * @param  {string} type
     *         The name of the console method to use.
     *
     * @param  {Array} args
     *         The arguments to be passed to the matching console method.
     */

    var LogByTypeFactory = function LogByTypeFactory(name, log) {
        return function (type, level, args) {
            var lvl = log.levels[level];
            var lvlRegExp = new RegExp("^(" + lvl + ")$");

            if (type !== 'log') {
                // Add the type to the front of the message when it's not "log".
                args.unshift(type.toUpperCase() + ':');
            } // Add console prefix after adding to history.


            args.unshift(name + ':'); // Add a clone of the args at this point to history.

            if (history) {
                history.push([].concat(args)); // only store 1000 history entries

                var splice = history.length - 1000;
                history.splice(0, splice > 0 ? splice : 0);
            } // If there's no console then don't try to output messages, but they will
            // still be stored in history.


            if (!window$1.console) {
                return;
            } // Was setting these once outside of this function, but containing them
            // in the function makes it easier to test cases where console doesn't exist
            // when the module is executed.


            var fn = window$1.console[type];

            if (!fn && type === 'debug') {
                // Certain browsers don't have support for console.debug. For those, we
                // should default to the closest comparable log.
                fn = window$1.console.info || window$1.console.log;
            } // Bail out if there's no console or if this type is not allowed by the
            // current logging level.


            if (!fn || !lvl || !lvlRegExp.test(type)) {
                return;
            }

            fn[Array.isArray(args) ? 'apply' : 'call'](window$1.console, args);
        };
    };

    function createLogger(name) {
        // This is the private tracking variable for logging level.
        var level = 'info'; // the curried logByType bound to the specific log and history

        var logByType;
        /**
         * Logs plain debug messages. Similar to `console.log`.
         *
         * Due to [limitations](https://github.com/jsdoc3/jsdoc/issues/955#issuecomment-313829149)
         * of our JSDoc template, we cannot properly document this as both a function
         * and a namespace, so its function signature is documented here.
         *
         * #### Arguments
         * ##### *args
         * Mixed[]
         *
         * Any combination of values that could be passed to `console.log()`.
         *
         * #### Return Value
         *
         * `undefined`
         *
         * @namespace
         * @param    {Mixed[]} args
         *           One or more messages or objects that should be logged.
         */

        var log = function log() {
            for (var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++) {
                args[_key] = arguments[_key];
            }

            logByType('log', level, args);
        }; // This is the logByType helper that the logging methods below use


        logByType = LogByTypeFactory(name, log);
        /**
         * Create a new sublogger which chains the old name to the new name.
         *
         * For example, doing `videojs.log.createLogger('player')` and then using that logger will log the following:
         * ```js
         *  mylogger('foo');
         *  // > VIDEOJS: player: foo
         * ```
         *
         * @param {string} name
         *        The name to add call the new logger
         * @return {Object}
         */

        log.createLogger = function (subname) {
            return createLogger(name + ': ' + subname);
        };
        /**
         * Enumeration of available logging levels, where the keys are the level names
         * and the values are `|`-separated strings containing logging methods allowed
         * in that logging level. These strings are used to create a regular expression
         * matching the function name being called.
         *
         * Levels provided by Video.js are:
         *
         * - `off`: Matches no calls. Any value that can be cast to `false` will have
         *   this effect. The most restrictive.
         * - `all`: Matches only Video.js-provided functions (`debug`, `log`,
         *   `log.warn`, and `log.error`).
         * - `debug`: Matches `log.debug`, `log`, `log.warn`, and `log.error` calls.
         * - `info` (default): Matches `log`, `log.warn`, and `log.error` calls.
         * - `warn`: Matches `log.warn` and `log.error` calls.
         * - `error`: Matches only `log.error` calls.
         *
         * @type {Object}
         */


        log.levels = {
            all: 'debug|log|warn|error',
            off: '',
            debug: 'debug|log|warn|error',
            info: 'log|warn|error',
            warn: 'warn|error',
            error: 'error',
            DEFAULT: level
        };
        /**
         * Get or set the current logging level.
         *
         * If a string matching a key from {@link module:log.levels} is provided, acts
         * as a setter.
         *
         * @param  {string} [lvl]
         *         Pass a valid level to set a new logging level.
         *
         * @return {string}
         *         The current logging level.
         */

        log.level = function (lvl) {
            if (typeof lvl === 'string') {
                if (!log.levels.hasOwnProperty(lvl)) {
                    throw new Error("\"" + lvl + "\" in not a valid log level");
                }

                level = lvl;
            }

            return level;
        };
        /**
         * Returns an array containing everything that has been logged to the history.
         *
         * This array is a shallow clone of the internal history record. However, its
         * contents are _not_ cloned; so, mutating objects inside this array will
         * mutate them in history.
         *
         * @return {Array}
         */


        log.history = function () {
            return history ? [].concat(history) : [];
        };
        /**
         * Allows you to filter the history by the given logger name
         *
         * @param {string} fname
         *        The name to filter by
         *
         * @return {Array}
         *         The filtered list to return
         */


        log.history.filter = function (fname) {
            return (history || []).filter(function (historyItem) {
                // if the first item in each historyItem includes `fname`, then it's a match
                return new RegExp(".*" + fname + ".*").test(historyItem[0]);
            });
        };
        /**
         * Clears the internal history tracking, but does not prevent further history
         * tracking.
         */


        log.history.clear = function () {
            if (history) {
                history.length = 0;
            }
        };
        /**
         * Disable history tracking if it is currently enabled.
         */


        log.history.disable = function () {
            if (history !== null) {
                history.length = 0;
                history = null;
            }
        };
        /**
         * Enable history tracking if it is currently disabled.
         */


        log.history.enable = function () {
            if (history === null) {
                history = [];
            }
        };
        /**
         * Logs error messages. Similar to `console.error`.
         *
         * @param {Mixed[]} args
         *        One or more messages or objects that should be logged as an error
         */


        log.error = function () {
            for (var _len2 = arguments.length, args = new Array(_len2), _key2 = 0; _key2 < _len2; _key2++) {
                args[_key2] = arguments[_key2];
            }

            return logByType('error', level, args);
        };
        /**
         * Logs warning messages. Similar to `console.warn`.
         *
         * @param {Mixed[]} args
         *        One or more messages or objects that should be logged as a warning.
         */


        log.warn = function () {
            for (var _len3 = arguments.length, args = new Array(_len3), _key3 = 0; _key3 < _len3; _key3++) {
                args[_key3] = arguments[_key3];
            }

            return logByType('warn', level, args);
        };
        /**
         * Logs debug messages. Similar to `console.debug`, but may also act as a comparable
         * log if `console.debug` is not available
         *
         * @param {Mixed[]} args
         *        One or more messages or objects that should be logged as debug.
         */


        log.debug = function () {
            for (var _len4 = arguments.length, args = new Array(_len4), _key4 = 0; _key4 < _len4; _key4++) {
                args[_key4] = arguments[_key4];
            }

            return logByType('debug', level, args);
        };

        return log;
    }

    /**
     * @file log.js
     * @module log
     */
    var log = createLogger('VIDEOJS');
    var createLogger$1 = log.createLogger;

    function createCommonjsModule(fn, module) {
        return module = { exports: {} }, fn(module, module.exports), module.exports;
    }

    var _extends_1 = createCommonjsModule(function (module) {
        function _extends() {
            module.exports = _extends = Object.assign || function (target) {
                for (var i = 1; i < arguments.length; i++) {
                    var source = arguments[i];

                    for (var key in source) {
                        if (Object.prototype.hasOwnProperty.call(source, key)) {
                            target[key] = source[key];
                        }
                    }
                }

                return target;
            };

            return _extends.apply(this, arguments);
        }

        module.exports = _extends;
    });

    /**
     * @file obj.js
     * @module obj
     */

    /**
     * @callback obj:EachCallback
     *
     * @param {Mixed} value
     *        The current key for the object that is being iterated over.
     *
     * @param {string} key
     *        The current key-value for object that is being iterated over
     */

    /**
     * @callback obj:ReduceCallback
     *
     * @param {Mixed} accum
     *        The value that is accumulating over the reduce loop.
     *
     * @param {Mixed} value
     *        The current key for the object that is being iterated over.
     *
     * @param {string} key
     *        The current key-value for object that is being iterated over
     *
     * @return {Mixed}
     *         The new accumulated value.
     */
    var toString = Object.prototype.toString;
    /**
     * Get the keys of an Object
     *
     * @param {Object}
     *        The Object to get the keys from
     *
     * @return {string[]}
     *         An array of the keys from the object. Returns an empty array if the
     *         object passed in was invalid or had no keys.
     *
     * @private
     */

    var keys = function keys(object) {
        return isObject(object) ? Object.keys(object) : [];
    };
    /**
     * Array-like iteration for objects.
     *
     * @param {Object} object
     *        The object to iterate over
     *
     * @param {obj:EachCallback} fn
     *        The callback function which is called for each key in the object.
     */


    function each(object, fn) {
        keys(object).forEach(function (key) {
            return fn(object[key], key);
        });
    }
    /**
     * Array-like reduce for objects.
     *
     * @param {Object} object
     *        The Object that you want to reduce.
     *
     * @param {Function} fn
     *         A callback function which is called for each key in the object. It
     *         receives the accumulated value and the per-iteration value and key
     *         as arguments.
     *
     * @param {Mixed} [initial = 0]
     *        Starting value
     *
     * @return {Mixed}
     *         The final accumulated value.
     */

    function reduce(object, fn, initial) {
        if (initial === void 0) {
            initial = 0;
        }

        return keys(object).reduce(function (accum, key) {
            return fn(accum, object[key], key);
        }, initial);
    }
    /**
     * Object.assign-style object shallow merge/extend.
     *
     * @param  {Object} target
     * @param  {Object} ...sources
     * @return {Object}
     */

    function assign(target) {
        for (var _len = arguments.length, sources = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
            sources[_key - 1] = arguments[_key];
        }

        if (Object.assign) {
            return _extends_1.apply(void 0, [target].concat(sources));
        }

        sources.forEach(function (source) {
            if (!source) {
                return;
            }

            each(source, function (value, key) {
                target[key] = value;
            });
        });
        return target;
    }
    /**
     * Returns whether a value is an object of any kind - including DOM nodes,
     * arrays, regular expressions, etc. Not functions, though.
     *
     * This avoids the gotcha where using `typeof` on a `null` value
     * results in `'object'`.
     *
     * @param  {Object} value
     * @return {boolean}
     */

    function isObject(value) {
        return !!value && typeof value === 'object';
    }
    /**
     * Returns whether an object appears to be a "plain" object - that is, a
     * direct instance of `Object`.
     *
     * @param  {Object} value
     * @return {boolean}
     */

    function isPlain(value) {
        return isObject(value) && toString.call(value) === '[object Object]' && value.constructor === Object;
    }

    /**
     * @file computed-style.js
     * @module computed-style
     */
    /**
     * A safe getComputedStyle.
     *
     * This is needed because in Firefox, if the player is loaded in an iframe with
     * `display:none`, then `getComputedStyle` returns `null`, so, we do a
     * null-check to make sure that the player doesn't break in these cases.
     *
     * @function
     * @param    {Element} el
     *           The element you want the computed style of
     *
     * @param    {string} prop
     *           The property name you want
     *
     * @see      https://bugzilla.mozilla.org/show_bug.cgi?id=548397
     */

    function computedStyle(el, prop) {
        if (!el || !prop) {
            return '';
        }

        if (typeof window$1.getComputedStyle === 'function') {
            var computedStyleValue = window$1.getComputedStyle(el);
            return computedStyleValue ? computedStyleValue.getPropertyValue(prop) || computedStyleValue[prop] : '';
        }

        return '';
    }

    /**
     * @file dom.js
     * @module dom
     */
    /**
     * Detect if a value is a string with any non-whitespace characters.
     *
     * @private
     * @param  {string} str
     *         The string to check
     *
     * @return {boolean}
     *         Will be `true` if the string is non-blank, `false` otherwise.
     *
     */

    function isNonBlankString(str) {
        // we use str.trim as it will trim any whitespace characters
        // from the front or back of non-whitespace characters. aka
        // Any string that contains non-whitespace characters will
        // still contain them after `trim` but whitespace only strings
        // will have a length of 0, failing this check.
        return typeof str === 'string' && Boolean(str.trim());
    }
    /**
     * Throws an error if the passed string has whitespace. This is used by
     * class methods to be relatively consistent with the classList API.
     *
     * @private
     * @param  {string} str
     *         The string to check for whitespace.
     *
     * @throws {Error}
     *         Throws an error if there is whitespace in the string.
     */


    function throwIfWhitespace(str) {
        // str.indexOf instead of regex because str.indexOf is faster performance wise.
        if (str.indexOf(' ') >= 0) {
            throw new Error('class has illegal whitespace characters');
        }
    }
    /**
     * Produce a regular expression for matching a className within an elements className.
     *
     * @private
     * @param  {string} className
     *         The className to generate the RegExp for.
     *
     * @return {RegExp}
     *         The RegExp that will check for a specific `className` in an elements
     *         className.
     */


    function classRegExp(className) {
        return new RegExp('(^|\\s)' + className + '($|\\s)');
    }
    /**
     * Whether the current DOM interface appears to be real (i.e. not simulated).
     *
     * @return {boolean}
     *         Will be `true` if the DOM appears to be real, `false` otherwise.
     */


    function isReal() {
        // Both document and window will never be undefined thanks to `global`.
        return document === window$1.document;
    }
    /**
     * Determines, via duck typing, whether or not a value is a DOM element.
     *
     * @param  {Mixed} value
     *         The value to check.
     *
     * @return {boolean}
     *         Will be `true` if the value is a DOM element, `false` otherwise.
     */

    function isEl(value) {
        return isObject(value) && value.nodeType === 1;
    }
    /**
     * Determines if the current DOM is embedded in an iframe.
     *
     * @return {boolean}
     *         Will be `true` if the DOM is embedded in an iframe, `false`
     *         otherwise.
     */

    function isInFrame() {
        // We need a try/catch here because Safari will throw errors when attempting
        // to get either `parent` or `self`
        try {
            return window$1.parent !== window$1.self;
        } catch (x) {
            return true;
        }
    }
    /**
     * Creates functions to query the DOM using a given method.
     *
     * @private
     * @param   {string} method
     *          The method to create the query with.
     *
     * @return  {Function}
     *          The query method
     */

    function createQuerier(method) {
        return function (selector, context) {
            if (!isNonBlankString(selector)) {
                return document[method](null);
            }

            if (isNonBlankString(context)) {
                context = document.querySelector(context);
            }

            var ctx = isEl(context) ? context : document;
            return ctx[method] && ctx[method](selector);
        };
    }
    /**
     * Creates an element and applies properties, attributes, and inserts content.
     *
     * @param  {string} [tagName='div']
     *         Name of tag to be created.
     *
     * @param  {Object} [properties={}]
     *         Element properties to be applied.
     *
     * @param  {Object} [attributes={}]
     *         Element attributes to be applied.
     *
     * @param {module:dom~ContentDescriptor} content
     *        A content descriptor object.
     *
     * @return {Element}
     *         The element that was created.
     */


    function createEl(tagName, properties, attributes, content) {
        if (tagName === void 0) {
            tagName = 'div';
        }

        if (properties === void 0) {
            properties = {};
        }

        if (attributes === void 0) {
            attributes = {};
        }

        var el = document.createElement(tagName);
        Object.getOwnPropertyNames(properties).forEach(function (propName) {
            var val = properties[propName]; // See #2176
            // We originally were accepting both properties and attributes in the
            // same object, but that doesn't work so well.

            if (propName.indexOf('aria-') !== -1 || propName === 'role' || propName === 'type') {
                log.warn('Setting attributes in the second argument of createEl()\n' + 'has been deprecated. Use the third argument instead.\n' + ("createEl(type, properties, attributes). Attempting to set " + propName + " to " + val + "."));
                el.setAttribute(propName, val); // Handle textContent since it's not supported everywhere and we have a
                // method for it.
            } else if (propName === 'textContent') {
                textContent(el, val);
            } else if (el[propName] !== val) {
                el[propName] = val;
            }
        });
        Object.getOwnPropertyNames(attributes).forEach(function (attrName) {
            el.setAttribute(attrName, attributes[attrName]);
        });

        if (content) {
            appendContent(el, content);
        }

        return el;
    }
    /**
     * Injects text into an element, replacing any existing contents entirely.
     *
     * @param  {Element} el
     *         The element to add text content into
     *
     * @param  {string} text
     *         The text content to add.
     *
     * @return {Element}
     *         The element with added text content.
     */

    function textContent(el, text) {
        if (typeof el.textContent === 'undefined') {
            el.innerText = text;
        } else {
            el.textContent = text;
        }

        return el;
    }
    /**
     * Insert an element as the first child node of another
     *
     * @param {Element} child
     *        Element to insert
     *
     * @param {Element} parent
     *        Element to insert child into
     */

    function prependTo(child, parent) {
        if (parent.firstChild) {
            parent.insertBefore(child, parent.firstChild);
        } else {
            parent.appendChild(child);
        }
    }
    /**
     * Check if an element has a class name.
     *
     * @param  {Element} element
     *         Element to check
     *
     * @param  {string} classToCheck
     *         Class name to check for
     *
     * @return {boolean}
     *         Will be `true` if the element has a class, `false` otherwise.
     *
     * @throws {Error}
     *         Throws an error if `classToCheck` has white space.
     */

    function hasClass(element, classToCheck) {
        throwIfWhitespace(classToCheck);

        if (element.classList) {
            return element.classList.contains(classToCheck);
        }

        return classRegExp(classToCheck).test(element.className);
    }
    /**
     * Add a class name to an element.
     *
     * @param  {Element} element
     *         Element to add class name to.
     *
     * @param  {string} classToAdd
     *         Class name to add.
     *
     * @return {Element}
     *         The DOM element with the added class name.
     */

    function addClass(element, classToAdd) {
        if (element.classList) {
            element.classList.add(classToAdd); // Don't need to `throwIfWhitespace` here because `hasElClass` will do it
            // in the case of classList not being supported.
        } else if (!hasClass(element, classToAdd)) {
            element.className = (element.className + ' ' + classToAdd).trim();
        }

        return element;
    }
    /**
     * Remove a class name from an element.
     *
     * @param  {Element} element
     *         Element to remove a class name from.
     *
     * @param  {string} classToRemove
     *         Class name to remove
     *
     * @return {Element}
     *         The DOM element with class name removed.
     */

    function removeClass(element, classToRemove) {
        if (element.classList) {
            element.classList.remove(classToRemove);
        } else {
            throwIfWhitespace(classToRemove);
            element.className = element.className.split(/\s+/).filter(function (c) {
                return c !== classToRemove;
            }).join(' ');
        }

        return element;
    }
    /**
     * The callback definition for toggleClass.
     *
     * @callback module:dom~PredicateCallback
     * @param    {Element} element
     *           The DOM element of the Component.
     *
     * @param    {string} classToToggle
     *           The `className` that wants to be toggled
     *
     * @return   {boolean|undefined}
     *           If `true` is returned, the `classToToggle` will be added to the
     *           `element`. If `false`, the `classToToggle` will be removed from
     *           the `element`. If `undefined`, the callback will be ignored.
     */

    /**
     * Adds or removes a class name to/from an element depending on an optional
     * condition or the presence/absence of the class name.
     *
     * @param  {Element} element
     *         The element to toggle a class name on.
     *
     * @param  {string} classToToggle
     *         The class that should be toggled.
     *
     * @param  {boolean|module:dom~PredicateCallback} [predicate]
     *         See the return value for {@link module:dom~PredicateCallback}
     *
     * @return {Element}
     *         The element with a class that has been toggled.
     */

    function toggleClass(element, classToToggle, predicate) {
        // This CANNOT use `classList` internally because IE11 does not support the
        // second parameter to the `classList.toggle()` method! Which is fine because
        // `classList` will be used by the add/remove functions.
        var has = hasClass(element, classToToggle);

        if (typeof predicate === 'function') {
            predicate = predicate(element, classToToggle);
        }

        if (typeof predicate !== 'boolean') {
            predicate = !has;
        } // If the necessary class operation matches the current state of the
        // element, no action is required.


        if (predicate === has) {
            return;
        }

        if (predicate) {
            addClass(element, classToToggle);
        } else {
            removeClass(element, classToToggle);
        }

        return element;
    }
    /**
     * Apply attributes to an HTML element.
     *
     * @param {Element} el
     *        Element to add attributes to.
     *
     * @param {Object} [attributes]
     *        Attributes to be applied.
     */

    function setAttributes(el, attributes) {
        Object.getOwnPropertyNames(attributes).forEach(function (attrName) {
            var attrValue = attributes[attrName];

            if (attrValue === null || typeof attrValue === 'undefined' || attrValue === false) {
                el.removeAttribute(attrName);
            } else {
                el.setAttribute(attrName, attrValue === true ? '' : attrValue);
            }
        });
    }
    /**
     * Get an element's attribute values, as defined on the HTML tag.
     *
     * Attributes are not the same as properties. They're defined on the tag
     * or with setAttribute.
     *
     * @param  {Element} tag
     *         Element from which to get tag attributes.
     *
     * @return {Object}
     *         All attributes of the element. Boolean attributes will be `true` or
     *         `false`, others will be strings.
     */

    function getAttributes(tag) {
        var obj = {}; // known boolean attributes
        // we can check for matching boolean properties, but not all browsers
        // and not all tags know about these attributes, so, we still want to check them manually

        var knownBooleans = ',' + 'autoplay,controls,playsinline,loop,muted,default,defaultMuted' + ',';

        if (tag && tag.attributes && tag.attributes.length > 0) {
            var attrs = tag.attributes;

            for (var i = attrs.length - 1; i >= 0; i--) {
                var attrName = attrs[i].name;
                var attrVal = attrs[i].value; // check for known booleans
                // the matching element property will return a value for typeof

                if (typeof tag[attrName] === 'boolean' || knownBooleans.indexOf(',' + attrName + ',') !== -1) {
                    // the value of an included boolean attribute is typically an empty
                    // string ('') which would equal false if we just check for a false value.
                    // we also don't want support bad code like autoplay='false'
                    attrVal = attrVal !== null ? true : false;
                }

                obj[attrName] = attrVal;
            }
        }

        return obj;
    }
    /**
     * Get the value of an element's attribute.
     *
     * @param {Element} el
     *        A DOM element.
     *
     * @param {string} attribute
     *        Attribute to get the value of.
     *
     * @return {string}
     *         The value of the attribute.
     */

    function getAttribute(el, attribute) {
        return el.getAttribute(attribute);
    }
    /**
     * Set the value of an element's attribute.
     *
     * @param {Element} el
     *        A DOM element.
     *
     * @param {string} attribute
     *        Attribute to set.
     *
     * @param {string} value
     *        Value to set the attribute to.
     */

    function setAttribute(el, attribute, value) {
        el.setAttribute(attribute, value);
    }
    /**
     * Remove an element's attribute.
     *
     * @param {Element} el
     *        A DOM element.
     *
     * @param {string} attribute
     *        Attribute to remove.
     */

    function removeAttribute(el, attribute) {
        el.removeAttribute(attribute);
    }
    /**
     * Attempt to block the ability to select text.
     */

    function blockTextSelection() {
        document.body.focus();

        document.onselectstart = function () {
            return false;
        };
    }
    /**
     * Turn off text selection blocking.
     */

    function unblockTextSelection() {
        document.onselectstart = function () {
            return true;
        };
    }
    /**
     * Identical to the native `getBoundingClientRect` function, but ensures that
     * the method is supported at all (it is in all browsers we claim to support)
     * and that the element is in the DOM before continuing.
     *
     * This wrapper function also shims properties which are not provided by some
     * older browsers (namely, IE8).
     *
     * Additionally, some browsers do not support adding properties to a
     * `ClientRect`/`DOMRect` object; so, we shallow-copy it with the standard
     * properties (except `x` and `y` which are not widely supported). This helps
     * avoid implementations where keys are non-enumerable.
     *
     * @param  {Element} el
     *         Element whose `ClientRect` we want to calculate.
     *
     * @return {Object|undefined}
     *         Always returns a plain object - or `undefined` if it cannot.
     */

    function getBoundingClientRect(el) {
        if (el && el.getBoundingClientRect && el.parentNode) {
            var rect = el.getBoundingClientRect();
            var result = {};
            ['bottom', 'height', 'left', 'right', 'top', 'width'].forEach(function (k) {
                if (rect[k] !== undefined) {
                    result[k] = rect[k];
                }
            });

            if (!result.height) {
                result.height = parseFloat(computedStyle(el, 'height'));
            }

            if (!result.width) {
                result.width = parseFloat(computedStyle(el, 'width'));
            }

            return result;
        }
    }
    /**
     * Represents the position of a DOM element on the page.
     *
     * @typedef  {Object} module:dom~Position
     *
     * @property {number} left
     *           Pixels to the left.
     *
     * @property {number} top
     *           Pixels from the top.
     */

    /**
     * Get the position of an element in the DOM.
     *
     * Uses `getBoundingClientRect` technique from John Resig.
     *
     * @see http://ejohn.org/blog/getboundingclientrect-is-awesome/
     *
     * @param  {Element} el
     *         Element from which to get offset.
     *
     * @return {module:dom~Position}
     *         The position of the element that was passed in.
     */

    function findPosition(el) {
        var box;

        if (el.getBoundingClientRect && el.parentNode) {
            box = el.getBoundingClientRect();
        }

        if (!box) {
            return {
                left: 0,
                top: 0
            };
        }

        var docEl = document.documentElement;
        var body = document.body;
        var clientLeft = docEl.clientLeft || body.clientLeft || 0;
        var scrollLeft = window$1.pageXOffset || body.scrollLeft;
        var left = box.left + scrollLeft - clientLeft;
        var clientTop = docEl.clientTop || body.clientTop || 0;
        var scrollTop = window$1.pageYOffset || body.scrollTop;
        var top = box.top + scrollTop - clientTop; // Android sometimes returns slightly off decimal values, so need to round

        return {
            left: Math.round(left),
            top: Math.round(top)
        };
    }
    /**
     * Represents x and y coordinates for a DOM element or mouse pointer.
     *
     * @typedef  {Object} module:dom~Coordinates
     *
     * @property {number} x
     *           x coordinate in pixels
     *
     * @property {number} y
     *           y coordinate in pixels
     */

    /**
     * Get the pointer position within an element.
     *
     * The base on the coordinates are the bottom left of the element.
     *
     * @param  {Element} el
     *         Element on which to get the pointer position on.
     *
     * @param  {EventTarget~Event} event
     *         Event object.
     *
     * @return {module:dom~Coordinates}
     *         A coordinates object corresponding to the mouse position.
     *
     */

    function getPointerPosition(el, event) {
        var position = {};
        var box = findPosition(el);
        var boxW = el.offsetWidth;
        var boxH = el.offsetHeight;
        var boxY = box.top;
        var boxX = box.left;
        var pageY = event.pageY;
        var pageX = event.pageX;

        if (event.changedTouches) {
            pageX = event.changedTouches[0].pageX;
            pageY = event.changedTouches[0].pageY;
        }

        position.y = Math.max(0, Math.min(1, (boxY - pageY + boxH) / boxH));
        position.x = Math.max(0, Math.min(1, (pageX - boxX) / boxW));
        return position;
    }
    /**
     * Determines, via duck typing, whether or not a value is a text node.
     *
     * @param  {Mixed} value
     *         Check if this value is a text node.
     *
     * @return {boolean}
     *         Will be `true` if the value is a text node, `false` otherwise.
     */

    function isTextNode(value) {
        return isObject(value) && value.nodeType === 3;
    }
    /**
     * Empties the contents of an element.
     *
     * @param  {Element} el
     *         The element to empty children from
     *
     * @return {Element}
     *         The element with no children
     */

    function emptyEl(el) {
        while (el.firstChild) {
            el.removeChild(el.firstChild);
        }

        return el;
    }
    /**
     * This is a mixed value that describes content to be injected into the DOM
     * via some method. It can be of the following types:
     *
     * Type       | Description
     * -----------|-------------
     * `string`   | The value will be normalized into a text node.
     * `Element`  | The value will be accepted as-is.
     * `TextNode` | The value will be accepted as-is.
     * `Array`    | A one-dimensional array of strings, elements, text nodes, or functions. These functions should return a string, element, or text node (any other return value, like an array, will be ignored).
     * `Function` | A function, which is expected to return a string, element, text node, or array - any of the other possible values described above. This means that a content descriptor could be a function that returns an array of functions, but those second-level functions must return strings, elements, or text nodes.
     *
     * @typedef {string|Element|TextNode|Array|Function} module:dom~ContentDescriptor
     */

    /**
     * Normalizes content for eventual insertion into the DOM.
     *
     * This allows a wide range of content definition methods, but helps protect
     * from falling into the trap of simply writing to `innerHTML`, which could
     * be an XSS concern.
     *
     * The content for an element can be passed in multiple types and
     * combinations, whose behavior is as follows:
     *
     * @param {module:dom~ContentDescriptor} content
     *        A content descriptor value.
     *
     * @return {Array}
     *         All of the content that was passed in, normalized to an array of
     *         elements or text nodes.
     */

    function normalizeContent(content) {
        // First, invoke content if it is a function. If it produces an array,
        // that needs to happen before normalization.
        if (typeof content === 'function') {
            content = content();
        } // Next up, normalize to an array, so one or many items can be normalized,
        // filtered, and returned.


        return (Array.isArray(content) ? content : [content]).map(function (value) {
            // First, invoke value if it is a function to produce a new value,
            // which will be subsequently normalized to a Node of some kind.
            if (typeof value === 'function') {
                value = value();
            }

            if (isEl(value) || isTextNode(value)) {
                return value;
            }

            if (typeof value === 'string' && /\S/.test(value)) {
                return document.createTextNode(value);
            }
        }).filter(function (value) {
            return value;
        });
    }
    /**
     * Normalizes and appends content to an element.
     *
     * @param  {Element} el
     *         Element to append normalized content to.
     *
     * @param {module:dom~ContentDescriptor} content
     *        A content descriptor value.
     *
     * @return {Element}
     *         The element with appended normalized content.
     */

    function appendContent(el, content) {
        normalizeContent(content).forEach(function (node) {
            return el.appendChild(node);
        });
        return el;
    }
    /**
     * Normalizes and inserts content into an element; this is identical to
     * `appendContent()`, except it empties the element first.
     *
     * @param {Element} el
     *        Element to insert normalized content into.
     *
     * @param {module:dom~ContentDescriptor} content
     *        A content descriptor value.
     *
     * @return {Element}
     *         The element with inserted normalized content.
     */

    function insertContent(el, content) {
        return appendContent(emptyEl(el), content);
    }
    /**
     * Check if an event was a single left click.
     *
     * @param  {EventTarget~Event} event
     *         Event object.
     *
     * @return {boolean}
     *         Will be `true` if a single left click, `false` otherwise.
     */

    function isSingleLeftClick(event) {
        // Note: if you create something draggable, be sure to
        // call it on both `mousedown` and `mousemove` event,
        // otherwise `mousedown` should be enough for a button
        if (event.button === undefined && event.buttons === undefined) {
            // Why do we need `buttons` ?
            // Because, middle mouse sometimes have this:
            // e.button === 0 and e.buttons === 4
            // Furthermore, we want to prevent combination click, something like
            // HOLD middlemouse then left click, that would be
            // e.button === 0, e.buttons === 5
            // just `button` is not gonna work
            // Alright, then what this block does ?
            // this is for chrome `simulate mobile devices`
            // I want to support this as well
            return true;
        }

        if (event.button === 0 && event.buttons === undefined) {
            // Touch screen, sometimes on some specific device, `buttons`
            // doesn't have anything (safari on ios, blackberry...)
            return true;
        } // `mouseup` event on a single left click has
        // `button` and `buttons` equal to 0


        if (event.type === 'mouseup' && event.button === 0 && event.buttons === 0) {
            return true;
        }

        if (event.button !== 0 || event.buttons !== 1) {
            // This is the reason we have those if else block above
            // if any special case we can catch and let it slide
            // we do it above, when get to here, this definitely
            // is-not-left-click
            return false;
        }

        return true;
    }
/**
 * Finds a single DOM element matching `selector` within the optional
 * `context` of another DOM element (defaulting to `document`).
 *
 * @param  {string} selector
 *         A valid CSS selector, which will be passed to `querySelector`.
 *
 * @param  {Element|String} [context=document]
 *         A DOM element within which to query. Can also be a selector
 *         string in which case the first matching element will be used