export default function Button({ onClick, text, className, ref }) {

    return (
        <button ref={ref} onClick={onClick} className={`bg-emerald-600 text-white rounded-md px-4 py-3 text-sm uppercase font-bold ${className}`}>
            {text}
        </button>
    )
}