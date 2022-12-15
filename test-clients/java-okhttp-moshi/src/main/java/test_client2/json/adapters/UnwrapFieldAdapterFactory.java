package test_client2.json.adapters;


import com.squareup.moshi.*;
import org.jetbrains.annotations.Nullable;

import java.io.IOException;
import java.lang.annotation.Annotation;
import java.lang.reflect.*;
import java.util.Set;

public final class UnwrapFieldAdapterFactory<T> implements JsonAdapter.Factory {
    final Class<T> type;

    public UnwrapFieldAdapterFactory(Class<T> type) {
        this.type = type;
    }

    @Nullable
    @Override
    public JsonAdapter<?> create(Type type, Set<? extends Annotation> annotations, Moshi moshi) {
        if (Types.getRawType(type) != this.type || !annotations.isEmpty()) {
            return null;
        }

        final Field[] fields = this.type.getDeclaredFields();
        if (fields.length != 1) {
            throw new RuntimeException("Type "+type.getTypeName()+" has "+fields.length+" fields, unwrap adapter can be used only with single-field classes");
        }
        var field = fields[0];

        var fieldName = field.getName();
        var getterName = "get" + fieldName.substring(0,1).toUpperCase() + fieldName.substring(1).toLowerCase();

        Method getter;
        try {
            getter = this.type.getDeclaredMethod(getterName);
        } catch (NoSuchMethodException e) {
            throw new RuntimeException("Type "+type.getTypeName()+" field "+fieldName+" does not have getter method "+ field.getType().getName()+", it's required for unwrap adapter", e);
        }

        Constructor<T> constructor;
        try {
            constructor = this.type.getDeclaredConstructor(field.getType());
        } catch (NoSuchMethodException e) {
            throw new RuntimeException("Type "+type.getTypeName()+" does not have constructor with single parameter of type "+ field.getType().getName()+", it's required for unwrap adapter", e);
        }

        JsonAdapter<?> fieldAdapter = moshi.adapter(field.getType());

        return new UnwrapFieldAdapter(constructor, getter, fieldAdapter);
    }

    public class UnwrapFieldAdapter<O, I> extends JsonAdapter<Object> {
        private Constructor<O> constructor;
        private Method getter;
        private JsonAdapter<I> fieldAdapter;

        public UnwrapFieldAdapter(Constructor<O> constructor, Method getter, JsonAdapter<I> fieldAdapter) {
            this.constructor = constructor;
            this.getter = getter;
            this.fieldAdapter = fieldAdapter;
        }

        @Override
        public Object fromJson(JsonReader reader) throws IOException {
            I fieldValue = fieldAdapter.fromJson(reader);
            try {
                return constructor.newInstance(fieldValue);
            } catch (Throwable e) {
                throw new IOException("Failed to create object with constructor "+constructor.getName(), e);
            }
        }

        @Override
        public void toJson(JsonWriter writer, Object value) throws IOException {
            if (value == null) {
                fieldAdapter.toJson(writer, null);
            } else {
                I fieldValue;
                try {
                    fieldValue = (I) getter.invoke(value);
                } catch (IllegalAccessException | InvocationTargetException e) {
                    throw new IOException("Failed to call method "+ getter.getName(), e);
                }
                fieldAdapter.toJson(writer, fieldValue);
            }
        }

        @Override
        public String toString() {
            return "UnwrapFieldAdapter(" + getter.getName() + ")";
        }
    }
}
