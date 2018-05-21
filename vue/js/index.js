// TODO: collision detection, autoscroll
// new: tab and highlighting, textareas, follow caret
var Autocomplete = Vue.component("autocomplete", {
	template: "#autocomplete-tpl",
	props: ["items", "placeholder", "label", "textarea", "rows", "cols"],
	data: function data() {
		return {
			id: 'input-' + parseInt(Math.random() * 1000),
			inputValue: "",
			searchMatch: [],
			selectedIndex: 0,
			clickedChooseItem: false,
			wordIndex: 0
		};
	},
	mounted: function mounted() {
		var _self = this;
		document.querySelector('#' + this.id).addEventListener('input', function () {
			var caret = getCaretCoordinates(this, this.selectionEnd);

			if (_self.searchMatch.length > 0) {
				var element = document.querySelectorAll('.' + _self.id + '-list');

				if (element[0]) {
					element[0].style.top = caret.top + 40 + 'px';
					element[0].style.left = caret.left + 'px';
				}
			}
		});
	},

	computed: {
		currentWord: function currentWord() {
			return this.inputValue.replace(/(\r\n|\n|\r)/gm, ' ').split(' ')[this.wordIndex];
		},
		inputSplitted: function inputSplitted() {
			return this.inputValue.replace(/(\r\n|\n|\r)/gm, ' ').split(" ");
		}
	},
	watch: {
		inputValue: function inputValue() {
			if (this.inputSplitted[this.inputSplitted.length-1] != "") {
				this.getSuggestions(this.inputSplitted[this.inputSplitted.length-1])
			}
			this.selectedIndex = 0;
			this.wordIndex = this.inputSplitted.length - 1;
		}
	},
	methods: {
		highlightWord: function highlightWord(word) {
			var regex = new RegExp("(" + this.currentWord + ")", "g");
			return word.replace(regex, '<mark>$1</mark>');
		},
		setWord: function setWord(word) {
			var currentWords = this.inputValue.replace(/(\r\n|\n|\r)/gm, '__br__ ').split(' ');
			currentWords[this.wordIndex] = currentWords[this.wordIndex].replace(this.currentWord, word + ' ');
			this.wordIndex += 1;
			this.inputValue = currentWords.join(' ').replace(/__br__\s/g, '\n');
		},
		moveDown: function moveDown() {
			if (this.selectedIndex < this.searchMatch.length - 1) {
				this.selectedIndex++;
			}
		},
		moveUp: function moveUp() {
			if (this.selectedIndex !== -1) {
				this.selectedIndex--;
			}
		},
		selectItem: function selectItem(index) {
			this.selectedIndex = index;
			this.chooseItem();
		},
		chooseItem: function chooseItem(e) {
			this.clickedChooseItem = true;

			if (this.selectedIndex !== -1 && this.searchMatch.length > 0) {
				if (e) {
					e.preventDefault();
				}
				this.setWord(this.searchMatch[this.selectedIndex]);
				this.selectedIndex = -1;
			}
		},
		focusout: function focusout(e) {
			var _this = this;

			setTimeout(function () {
				if (!_this.clickedChooseItem) {
					_this.searchMatch = [];
					_this.selectedIndex = -1;
				}
				_this.clickedChooseItem = false;
			}, 100);
		},
		focus: function focus() {
			var _this2 = this;

			this.searchMatch = [];
			if (this.currentWord !== "" && this.items != undefined) {
				this.searchMatch = this.items.filter(function (el) {
					return el.indexOf(_this2.currentWord) >= 0;
				});
			}
			if (this.searchMatch.length === 1 && this.currentWord === this.searchMatch[0]) {
				this.searchMatch = [];
			}
		},
		getSuggestions: function getSuggestions(word) {
			this.$http.get('/suggest/'+word, function(data, status, request){
				if(status == 200)
				{
				  this.items = data;
				  this.focus();
				}
			})
		}
	}
});

new Vue({
	el: "#app"
});

// Thanks: https://github.com/component/textarea-caret-position
(function () {
	function e(b, e, f) {
		if (!h) throw Error("textarea-caret-position#getCaretCoordinates should only be called in a browser");if (f = f && f.debug || !1) {
			var a = document.querySelector("#input-textarea-caret-position-mirror-div");a && a.parentNode.removeChild(a);
		}a = document.createElement("div");a.id = "input-textarea-caret-position-mirror-div";document.body.appendChild(a);var c = a.style,
		    d = window.getComputedStyle ? window.getComputedStyle(b) : b.currentStyle,
		    k = "INPUT" === b.nodeName;c.whiteSpace = "pre-wrap";k || (c.wordWrap = "break-word");c.position = "absolute";f || (c.visibility = "hidden");l.forEach(function (a) {
			k && "lineHeight" === a ? c.lineHeight = d.height : c[a] = d[a];
		});m ? b.scrollHeight > parseInt(d.height) && (c.overflowY = "scroll") : c.overflow = "hidden";a.textContent = b.value.substring(0, e);k && (a.textContent = a.textContent.replace(/\s/g, "\xA0"));var g = document.createElement("span");g.textContent = b.value.substring(e) || ".";a.appendChild(g);b = { top: g.offsetTop + parseInt(d.borderTopWidth), left: g.offsetLeft + parseInt(d.borderLeftWidth), height: parseInt(d.lineHeight) };
		f ? g.style.backgroundColor = "#aaa" : document.body.removeChild(a);return b;
	}var l = "direction boxSizing width height overflowX overflowY borderTopWidth borderRightWidth borderBottomWidth borderLeftWidth borderStyle paddingTop paddingRight paddingBottom paddingLeft fontStyle fontVariant fontWeight fontStretch fontSize fontSizeAdjust lineHeight fontFamily textAlign textTransform textIndent textDecoration letterSpacing wordSpacing tabSize MozTabSize".split(" "),
	    h = "undefined" !== typeof window,
	    m = h && null != window.mozInnerScreenX;
	"undefined" != typeof module && "undefined" != typeof module.exports ? module.exports = e : h && (window.getCaretCoordinates = e);
})();