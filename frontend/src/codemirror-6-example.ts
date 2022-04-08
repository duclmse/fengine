import { EditorView, ViewPlugin, ViewUpdate } from "@codemirror/next/view";
import { EditorState, StateField } from "@codemirror/next/state";
// import { defaultKeymap } from "@codemirror/next/commands";
import { basicSetup } from "@codemirror/next/basic-setup";
import { javascript } from "@codemirror/next/lang-javascript";

export const run = () => {
    const app = document.getElementById("app");

    let countDocChanges = StateField.define<number>({
        create() {
            return 0;
        },
        update(value, tr) {
            return tr.docChanged ? value + 1 : value;
        }
    });

    const docSizePlugin = ViewPlugin.fromClass(
        class {
            private dom: HTMLDivElement;
            constructor(view: EditorView) {
                this.dom = view.dom.appendChild(document.createElement("div"));
                this.dom.style.cssText =
                    'color: purple; padding: 1rem; border: 1px solid gray';
                this.dom.textContent =
                    `doc length: ${view.state.doc.length}
 doc changes: ${view.state.field(countDocChanges)}`;
            }

            update(update: ViewUpdate) {
                if (update.docChanged)
                    this.dom.textContent =
                        `doc length: ${update.state.doc.length}
 doc changes: ${update.state.field(countDocChanges)}`;
            }

            destroy() {
                this.dom.remove();
            }
        }
    );

    let state = EditorState.create({
        doc: "const x = () => 'hi!'",
        extensions: [basicSetup, javascript(), countDocChanges, docSizePlugin]
    });

    let view = new EditorView({
        state,
        parent: app!
    });
};
