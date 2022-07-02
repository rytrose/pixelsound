const Button = ({ children, onClick }) => {
  return (
    <button
      className="outline-0 select-none p-2 rounded-xl text-black bg-violet-300 hover:bg-violet-400 active:bg-violet-500"
      onClick={onClick}
    >
      {children}
    </button>
  );
};

export default Button;
