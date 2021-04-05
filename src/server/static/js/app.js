class TranslationRequest {
  _memory = {};

  // cb: ({ translation: string }) => void
  constructor(cb) {
    this._cb = cb;
  }

  async translate(
    options = { sourceLangauge, targetLanguage, sentence }
  ) {
    options = {
      sentence: "",
      sourceLangauge: "auto",
      targetLanguage: "en",

      ...options,
    };

    if (
      this._memory[options.sentence] &&
      this._memory[options.sentence][options.targetLanguage]
    )
      return;
    this._memory[options.sentence] ||= {};
    this._memory[options.sentence][options.targetLanguage] = true;

    try {
      console.log(`fetching translation for ${options.sentence}`);
      const res = await fetch(
        `api/translate?sentence=${options.sentence}&to=${options.targetLanguage}&from=${options.sourceLangauge}`
      );
      const resJSON = await res.json();
      this._cb(resJSON);
    } catch (e) {
      console.log(e);
    } finally {
      delete this._memory[options.sentence][options.targetLanguage];
      if (
        Object.getOwnPropertyNames(this._memory[options.sentence])
          .length === 0
      ) {
        delete this._memory[options.sentence];
      }
    }
  }
}

class Emitter {
  _callbacks = [];
  onUpdate(cb) {
    this._callbacks.push(cb);
  }
  emit() {
    this._callbacks.forEach((cb) => cb());
  }
}

const modelUpdatedEmitter = new Emitter();

const textHistory = {
  // { sentence, translation }
  queue: [],
  max: 99,
  current: 0,

  isEmpty() {
    return this.queue.length === 0;
  },
  isLast() {
    if (this.queue.length === 0) return true;
    return this.current === this.queue.length - 1;
  },
  isFirst() {
    return this.current === 0;
  },

  fixCurrentBounded() {
    if (this.queue.length === 0) {
      this.current = 0;
    } else {
      if (this.current >= this.queue.length) {
        this.last();
      }
    }
    if (this.current < 0) this.first();
    modelUpdatedEmitter.emit();
  },

  push(sentence) {
    this.queue.push({ sentence });
    while (this.queue.length > this.max) {
      this.queue.shift();
    }
    this.last();
    modelUpdatedEmitter.emit();
  },

  addTranslation(options = { sentence, translation, targetLanguage }) {
    const { sentence, translation, targetLanguage } = options;
    console.log("addTranslation:", options);
    this.queue = this.queue.map((el, idx) => {
      if (el.sentence !== sentence) return el;
      el.translation ||= {};
      el.translation[targetLanguage] = translation;
      return el;
    });
    modelUpdatedEmitter.emit();
  },

  removeTranslation(options = { sentence, targetLanguage }) {
    console.log("removeTranslation", options);
    const { sentence, targetLanguage } = options;
    this.queue = this.queue.map((el, idx) => {
      if (el.sentence !== sentence || !el.translation) return el;
      delete el.translation[targetLanguage];
      return el;
    });
    modelUpdatedEmitter.emit();
  },

  delete() {
    if (this.queue.length === 0) return;
    this.queue.splice(this.current, 1);
    this.prev();
    this.fixCurrentBounded();
    modelUpdatedEmitter.emit();
  },
  prev() {
    this.current--;
    this.fixCurrentBounded();
    modelUpdatedEmitter.emit();
  },
  next() {
    this.current++;
    this.fixCurrentBounded();
    modelUpdatedEmitter.emit();
  },
  first() {
    this.current = 0;
    modelUpdatedEmitter.emit();
  },
  last() {
    this.current = this.queue.length === 0 ? 0 : this.queue.length - 1;
    modelUpdatedEmitter.emit();
  },

  getCurrent() {
    if (this.current < 0) return;
    return this.queue[this.current];
  },
};

const clipboard = {
  _lastTextLine: "",

  setLastTextLine(text) {
    clipboard._lastTextLine = text;
  },

  copy(containerid) {
    const result = {
      copied: true,
    };

    const element = document.getElementById(containerid);
    if (element.textContent === this._lastTextLine)
      return { ...result, copied: false };

    const range = document.createRange();
    range.selectNode(document.getElementById(containerid));
    window.getSelection().addRange(range);
    document.execCommand("copy");
    window.getSelection().removeRange(range);
    return result;
  },
};

const clipboardObserver = {
  observer: new MutationObserver(() => {
    let ps = document.getElementsByTagName("p");

    // removing all except last
    while (ps.length > 0) {
      const sentence = ps[0].textContent;
      ps[0].remove();
    }
  }),

  start() {
    this.observer.observe(document.getElementById("body"), {
      childList: true,
    });
  },
};
_.extend(clipboardObserver, Backbone.Events);

const gui = {
  init() {
    this.buttons.init();
    this.langs.init();
    this.update();
    modelUpdatedEmitter.onUpdate(this.update.bind(this));
  },

  langs: {
    elements: Array.from(document.getElementById("langs").children),
    init() {
      this.elements.forEach((el) => {
        el.onchange = function () {
          modelUpdatedEmitter.emit();
        };
      });
    },
    getActiveLangs() {
      return this.elements
        .filter((el) => el.checked)
        .map((el) => el.name);
    },
  },

  buttons: {
    copyButton: document.getElementById("copy_button"),
    deleteButton: document.getElementById("delete_button"),
    prevButton: document.getElementById("prev_button"),
    nextButton: document.getElementById("next_button"),
    lastButton: document.getElementById("last_button"),

    updatePrevButton() {
      if (textHistory.isFirst()) {
        this.prevButton.disabled = true;
      } else {
        this.prevButton.disabled = false;
      }
    },
    updateNextButton() {
      if (textHistory.isLast()) {
        this.nextButton.disabled = true;
      } else {
        this.nextButton.disabled = false;
      }
    },
    updateDeleteButton() {
      if (textHistory.isEmpty()) {
        this.deleteButton.disabled = true;
      } else {
        this.deleteButton.disabled = false;
      }
    },

    init() {
      this.copyButton.onclick = () => {
        const { copied } = clipboard.copy("current_line");
        clipboardObserver.skipNextClipboard |= copied;
      };

      this.prevButton.onclick = () => {
        textHistory.prev();
      };

      this.nextButton.onclick = () => {
        textHistory.next();
      };

      this.deleteButton.onclick = () => {
        textHistory.delete();
      };

      this.lastButton.onclick = () => {
        textHistory.last();
      };
    },
  },

  update() {
    this.showCurrentLine();
    this.showTextHistoryStatus();
    this.updateButtons();
  },
  showCurrentLine() {
    const current = textHistory.getCurrent();
    if (!current) return;
    const { translation, sentence } = current;

    let currentLineElement = document.getElementById("current_line");
    currentLineElement.textContent = sentence;

    let anchor = document.getElementById("translation");
    anchor.innerHTML = "";
    const targetLanguages = gui.langs.getActiveLangs();
    targetLanguages.forEach((targetLanguage) => {
      if (translation && translation[targetLanguage]) {
        const translationDisplayElement = document.createElement("div");
        translationDisplayElement.onclick = () => {
          console.log(`${targetLanguage} clicked`);
          textHistory.removeTranslation({ sentence, targetLanguage });
        };
        translationDisplayElement.setAttribute(
          "class",
          `lang-${targetLanguage}`
        );
        translationDisplayElement.textContent =
          translation[targetLanguage];
        anchor.appendChild(translationDisplayElement);
      }
    });
  },
  showTextHistoryStatus() {
    let textHistoryStatusElement = document.getElementById(
      "text_history_status"
    );

    function makeTextContent(current, max) {
      function toString(num) {
        return num.toString().padStart(2, "0");
      }

      return `${toString(current)}/${toString(max)}`;
    }

    textHistoryStatusElement.textContent = textHistory.isEmpty()
      ? makeTextContent(0, 0)
      : makeTextContent(
        textHistory.current + 1,
        textHistory.queue.length
      );
  },
  updateButtons() {
    this.buttons.updateDeleteButton();
    this.buttons.updateNextButton();
    this.buttons.updatePrevButton();
  },
};

clipboardObserver.start();
gui.init();

const translationElement = document.getElementById("translation");
const translationRequest = new TranslationRequest(
  ({ translation, sentence, targetLanguage }) => {
    textHistory.addTranslation({ sentence, translation, targetLanguage });
  }
);

modelUpdatedEmitter.onUpdate(() => {
  const current = textHistory.getCurrent();
  if (!current) return;
  const { sentence, translation } = current;

  const targetLanguages = gui.langs.getActiveLangs();
  targetLanguages.forEach((targetLanguage) => {
    if (translation && translation[targetLanguage]) return;

    translationRequest.translate({
      sentence,
      targetLanguage,
    });
  });
});