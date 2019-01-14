/* globals window */
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

				if (el.addEventListener) {
					el.addEventListener('click', dblClick, false);
				} else if (el.attachEvent) {
					el.attachEvent('onclick', dblClick);
				}
			},
			'autoselect': function (el) {
				setTimeout(function () {
					el.select();
				});
			},
			'blur-submit': function (el) {
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
				if (el.addEventListener) {
					el.addEventListener('blur', onBlur, false);
				} else if (el.attachEvent) {
					el.attachEvent('onblur', onBlur);
				}
			}
		}
	});

})(window);
