const FileInput = ({ accept = undefined, children }) => {
  return (
    <div className="py-2">
      <label className="outline-none focus-within:ring-2 focus-within:ring-violet-500 select-none cursor-pointer p-2 rounded-xl text-black bg-violet-300 hover:bg-violet-400 active:bg-violet-500">
        {children}
        <input
          className="opacity-0 absolute left-[-99999rem]"
          type="file"
          accept={accept}
        ></input>
      </label>
    </div>
  );
};

export default FileInput;
