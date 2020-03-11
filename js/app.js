/* globals window */
/*
Cheap and cheerful JS event handlers used to bind
user events for crafting treetop XHR requests.
*/

(function (window) {
	'use strict';

	window.treetop.init({
		'mountAttrs': {
			'input-submit': (elm) => {
				// allows a single input element to be treated as if it is a form.
				// for the purpose of treetop requests
				const url = elm.getAttribute('input-submit');
				if (url === null) {
					return;
				}
				var type = (elm.getAttribute('type') || 'text').toUpperCase();
				var method = elm.getAttribute('method') || 'GET';
				switch (method.toUpperCase()) {
				case 'POST':
					method = 'POST';
					break;

				case 'GET':
				case '':
					method = 'GET';
					break;

				default:
					throw new Error('input-submit component: unsupported method ' + method);
				}

				elm.addEventListener('click', () => {
					const data = new window.FormData();
					if ((type != 'CHECKBOX' && type != 'RADIO') || elm.checked) {
						data.append(elm.name, elm.value);
					}
					const body = new window.URLSearchParams(data).toString();
					window.treetop.request(method, url, body, 'application/x-www-form-urlencoded');
				});
			},
			'double-click-link': function (el) {
				// equivalent to treetop-link except activating it requires
				// two successive clicks within 0.8 seconds.
				function dblClick(_evt) {
					var evt = _evt || window.event;
					var elm = evt.target || evt.srcElement;
					if (!elm.hasAttribute('double-click-link')) {
						return;
					}
					var ts = Date.now();
					if (elm.__lastClicked__ && ts - elm.__lastClicked__ < 800 ) {
						var href = elm.getAttribute('double-click-link');
						window.treetop.request('GET', href);
						elm.__lastClicked__ = null;
					} else {
						elm.__lastClicked__ = ts;
					}
				}
				el.addEventListener('click', dblClick, false);
			},
			'autoselect': function (el) {
				// select the contents of an input element as soon as it is mounted to the DOM
				setTimeout(function () {
					el.select();
				});
			},
			'blur-submit': function (el) {
				// blur even on an input element will cause the enclosing form to be submitted
				// TODO: This is a crude implementation, combined with input-submit it is causing a double-tap currently
				function onBlur(_evt) {
					var evt = _evt || window.event;
					var elm = evt.target || evt.srcElement;
					while (elm.tagName.toUpperCase() !== 'FORM') {
						if (elm.parentElement) {
							elm = elm.parentElement;
						} else {
							return; // this is not wrapped in a form
						}
					}
					if (elm.hasAttribute('treetop')) {
						window.treetop.submit(elm);
					} else {
						elm.submit();
					}
				}
				el.addEventListener('blur', onBlur, false);
			}
		}
	});

})(window);
