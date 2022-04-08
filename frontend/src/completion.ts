import {autocompletion, CompletionContext} from "@codemirror/autocomplete"
import {basicSetup} from "@codemirror/next/basic-setup";

function myCompletions(context: CompletionContext) {
    let word = context.matchBefore(/\w*/);
    if (word === null || word.from == word.to && !context.explicit) return null;

    return {
        from: word.from,
        options: [
            {label: "match", type: "keyword"},
            {label: "hello", type: "variable", info: "(World)"},
            {label: "magic", type: "text", apply: "⠁⭒*.✩.*⭒⠁", detail: "macro"}
        ]
    }
}

let state = EditorState.create({
    doc: "Press Ctrl-Space in here...\n",
    extensions: [basicSetup, autocompletion({override: [myCompletions]})]
})
