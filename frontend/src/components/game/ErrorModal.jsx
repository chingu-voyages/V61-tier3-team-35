export default function ErrorModal({ error }) {
    return (
        <article role="alert" className="absolute z-20 bg-white top-1/4 rounded-xl overflow-hidden shadow-2xl px-16 py-8 text-center font-bold text-gray-800">
            {error}
        </article>
    )
}