import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router";
import { Link } from "react-router-dom";
import { marked } from "marked"
import { useSelector } from "react-redux";

function paginate(array, page_size, page_number) {
  return array.slice((page_number - 1) * page_size, page_number * page_size);
}

export default function List() {
    

    const url = useLocation()
    const type = url.pathname.split('/')[1]
    const locale = useSelector(state => state.locale.value)

    const [list, setList] = useState([])
    const [perpage, setPerpage] = useState(5)
    const [page, setPage] = useState(1)
    
    useEffect(() => {
    
        console.log(url, type, locale)
  
        const prefix = import.meta.env.VITE_BLOG_PREFIX;
        fetch(prefix + type + "/list_" + locale + ".json").then((response) => {
            response.json().then((json) => {
                setList(json)
            })
        })
        
    }, [locale, url, locale])

  return (
    <div>
      <h1>List</h1>
        <ul>
            {
                paginate(list, perpage, page).map((item, idx) => (
                    type === "posts" ? <li key={"post" + idx}>
                        <Link  to={"/posts/" + item.Folder + "/" + item.Id}>{item.Title}</Link>
                    </li>:
                    <li key={"restricted" + idx}>
                        <Link  to={"/restricted/" + item.Hash + "/" + item.Id}>{item.Title}</Link>
                    </li>
                ))
            }
        </ul>
        <ul>
            {/* page numbers */}
            {
                Array.from(Array(Math.ceil(list.length / perpage)).keys()).map((item, idx) => (
                    <li key={"page" + idx} >
                        <Link onClick={(e) =>{
                            e.preventDefault()
                            setPage(item + 1)
                        }}>{item + 1}</Link>
                    </li>
                ))
            }
        </ul>
    </div>
  );
}
