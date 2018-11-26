/* globals window */
(function (window) {
	'use strict';

	window.treetop.init({
		'mountAttrs': {
			'single-action': (elm) => {
				// allows a single input element to be treated as if it is a form.
				// for the purpose of treetop requests
				const url = elm.getAttribute('single-action');
				var type = (elm.getAttribute('type') || 'text').toUpperCase();
				var method = elm.getAttribute('method');
				switch (method.toUpperCase()) {
				case 'POST':
					method = 'POST';
					break;

				case 'GET':
				case '':
					method = 'GET';
					break;

				default:
					throw new Error('Single-action component: unsupported method ' + method);
				}

				elm.addEventListener('click', () => {
					const data = new window.FormData();
					if ((type != 'CHECKBOX' && type != 'RADIO') || elm.checked) {
						data.append(elm.name, elm.value);
					}
					const body = new window.URLSearchParams(data).toString();
					window.treetop.request(method, url, body, 'application/x-www-form-urlencoded');
				});
			}
		}
	});

})(window);
