const Header = ({ className }) => {
  return (
    <div className={"flex items-center justify-between px-3 py-2 " + className}>
      <p className="font-serif text-2xl text-black">pixelsound</p>
    </div>
  );
};

export default Header;
